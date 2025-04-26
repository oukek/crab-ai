package llm

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"oukek/crab-ai/app/common"
	"oukek/crab-ai/app/common/sse"
	"oukek/crab-ai/app/service/chains"
	_type "oukek/crab-ai/app/service/llm/type"

	"github.com/gin-gonic/gin"
)

// MakeChatWithStream creates a handler for streaming chat responses.
// It sets up the SSE headers for Gin context or accumulates the response if context is nil.
func MakeChatWithStream(c *gin.Context) chains.HandlerFunc[*ChatChains] {
	// --- Handler for non-HTTP context (e.g., internal calls) ---
	if c == nil {
		return func(chain *ChatChains) error {
			messageChan := make(chan _type.MessageChanItem, 100)
			// Run ChatWithStream in a goroutine
			go common.WithSafeFn(func() {
				// Pass context from chain
				if err := ChatWithStream(chain, messageChan); err != nil {
					// Send error back through channel if ChatWithStream fails early
					select {
					case messageChan <- _type.MessageChanItem{Event: "err", Data: err.Error()}:
					default:
						common.BaseLogger.WithError(err).Error("Failed to send initial error to messageChan")
					}
					close(messageChan) // Ensure channel is closed on error
				}
			})

			var contentBuilder strings.Builder
			finalResponse := _type.Response{}

			common.BaseLogger.Debug("Internal stream: Waiting for messages...")
			for msg := range messageChan {
				switch msg.Event {
				case "err":
					common.BaseLogger.Errorf("Internal stream: Received error: %v", msg.Data)
					// Attempt to convert error data to error type
					if errStr, ok := msg.Data.(string); ok {
						return errors.New(errStr)
					} else if err, ok := msg.Data.(error); ok {
						return err
					} else {
						return fmt.Errorf("received unknown error type: %v", msg.Data)
					}
				case "init":
					if item, ok := msg.Data.(_type.StreamItem); ok {
						common.BaseLogger.Debugf("Internal stream: Received init: %+v", item)
						finalResponse.ID = item.Id
						finalResponse.Object = strings.Replace(item.Object, ".chunk", "", 1) // Convert chunk object to completion object
						finalResponse.Model = item.Model
						finalResponse.Created = item.Created
					} else {
						common.BaseLogger.Warnf("Internal stream: Received init with unexpected data type: %T", msg.Data)
					}
				case "data":
					if strData, ok := msg.Data.(string); ok {
						common.BaseLogger.Debugf("Internal stream: Received data: %s", strData)
						contentBuilder.WriteString(strData)
					} else {
						common.BaseLogger.Warnf("Internal stream: Received data with unexpected type: %T", msg.Data)
					}
				case "finish":
					common.BaseLogger.Debug("Internal stream: Received finish event.")
					// Final assembly happens after loop
				}
			}

			common.BaseLogger.Debugf("Internal stream: Finished receiving. Final content length: %d", contentBuilder.Len())

			// Assemble final response
			finalResponse.Choices = []_type.NonStreamingChoice{
				{
					Index:        0,
					FinishReason: "stop", // Assuming stop, could be different based on last chunk
					Message: &_type.ResMessage{
						Role:    "assistant",
						Content: contentBuilder.String(),
					},
				},
			}
			// TODO: Aggregate usage if available in stream items
			// if firstItem != nil && firstItem.Usage != nil { finalResponse.Usage = firstItem.Usage }

			chain.Res = &finalResponse // Store the aggregated response in the chain
			return nil
		}
	}

	// --- Handler for Gin HTTP context ---
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Nginx: prevent buffering

	return func(chain *ChatChains) error {
		messageChan := make(chan _type.MessageChanItem, 100)
		clientGone := c.Writer.CloseNotify()
		requestContext := c.Request.Context() // Use Gin request context

		// Run ChatWithStream in a goroutine
		go common.WithSafeFn(func() {
			// Pass request context to ChatWithStream
			if err := ChatWithStream(chain, messageChan); err != nil {
				// Send error back through channel if ChatWithStream fails early
				select {
				case messageChan <- _type.MessageChanItem{Event: "err", Data: err.Error()}:
				case <-requestContext.Done(): // Avoid blocking if context is cancelled
				}
				close(messageChan) // Ensure channel is closed on error
			}
		})

		streamClosed := false              // Flag to prevent writing after error or client disconnect
		var contentBuilder strings.Builder // To store aggregated content for final response
		finalResponse := _type.Response{}

		common.BaseLogger.Debug("HTTP stream: Waiting for messages...")
		for {
			select {
			case <-clientGone:
				common.BaseLogger.Info("HTTP stream: Client disconnected.")
				streamClosed = true
				return nil // Or return an error indicating client disconnect
			case <-requestContext.Done():
				common.BaseLogger.Info("HTTP stream: Request context cancelled.")
				streamClosed = true
				return requestContext.Err() // Return context error
			case msg, ok := <-messageChan:
				if !ok {
					common.BaseLogger.Debug("HTTP stream: Message channel closed.")
					if !streamClosed {
						// Send DONE event if stream hasn't been closed by error/disconnect
						_, err := fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
						if err == nil {
							c.Writer.Flush()
						}
					}
					// Assemble final response after loop
					finalResponse.Choices = []_type.NonStreamingChoice{
						{
							Index:        0,
							FinishReason: "stop",
							Message: &_type.ResMessage{
								Role:    "assistant",
								Content: contentBuilder.String(),
							},
						},
					}
					chain.Res = &finalResponse
					return nil // Normal stream completion
				}

				if streamClosed {
					continue // Don't process messages if stream is already considered closed
				}

				switch msg.Event {
				case "err":
					common.BaseLogger.Errorf("HTTP stream: Received error: %v", msg.Data)
					// Send error event to client
					_, err := fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", msg.Data)
					if err == nil {
						c.Writer.Flush()
					}
					streamClosed = true // Stop sending further events
					// Return the error that occurred during streaming
					if errStr, ok := msg.Data.(string); ok {
						return errors.New(errStr)
					} else if err, ok := msg.Data.(error); ok {
						return err
					} else {
						return fmt.Errorf("received unknown error type: %v", msg.Data)
					}
				case "init":
					if item, ok := msg.Data.(_type.StreamItem); ok {
						common.BaseLogger.Debugf("HTTP stream: Received init: %+v", item)
						// Store initial data for final response assembly
						finalResponse.ID = item.Id
						finalResponse.Object = strings.Replace(item.Object, ".chunk", "", 1)
						finalResponse.Model = item.Model
						finalResponse.Created = item.Created
						// Do not send init event via SSE, only data chunks
					} else {
						common.BaseLogger.Warnf("HTTP stream: Received init with unexpected data type: %T", msg.Data)
					}
				case "data":
					if strData, ok := msg.Data.(string); ok {
						common.BaseLogger.Debugf("HTTP stream: Sending data chunk: %s", strData)
						contentBuilder.WriteString(strData) // Aggregate content
						// Format as SSE data event
						sseData := strings.ReplaceAll(strData, "\n", "\\n") // Escape newlines for SSE
						_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", sseData)
						if err != nil {
							common.BaseLogger.WithError(err).Error("HTTP stream: Failed to write data to client")
							streamClosed = true // Stop on write error
						} else {
							c.Writer.Flush() // Flush data immediately
						}
					} else {
						common.BaseLogger.Warnf("HTTP stream: Received data with unexpected type: %T", msg.Data)
					}
				case "finish":
					common.BaseLogger.Debug("HTTP stream: Received finish event. Channel will close.")
					// Don't send finish via SSE, client detects end when channel closes or gets [DONE]
				}
			}
		}
	}
}

// ChatWithStream performs the actual streaming request and processing.
// It uses the LLMProvider to handle API specifics and sends events via messageChan.
func ChatWithStream(chain *ChatChains, messageChan chan<- _type.MessageChanItem) (err error) {
	// Ensure messageChan is closed eventually
	defer func() {
		if r := recover(); r != nil {
			// Log panic and try to send error through channel
			emsg := fmt.Sprintf("Panic recovered in ChatWithStream: %v", r)
			common.BaseLogger.Error(emsg)
			select {
			case messageChan <- _type.MessageChanItem{Event: "err", Data: emsg}:
			default: // Avoid blocking if channel is full/closed
			}
			err = errors.New(emsg) // Set return error
		}
		close(messageChan)
		common.BaseLogger.Debug("ChatWithStream finished and channel closed.")
	}()

	common.BaseLogger.Debug("ChatWithStream started.")

	// 1. Create Request Parameters using Provider
	provider, err := GetProvider(chain.Provider.Type)
	if err != nil {
		common.BaseLogger.WithError(err).Error("Failed to get provider")
		messageChan <- _type.MessageChanItem{Event: "err", Data: fmt.Sprintf("Failed to get provider: %v", err)}
		return err // Return error early
	}

	reqParams, err := provider.MakeStreamChatRequest(context.Background(), &chain.Provider, chain.Req) // Assuming provider has access to its model.Provider info
	if err != nil {
		common.BaseLogger.WithError(err).Error("Failed to create stream request parameters")
		messageChan <- _type.MessageChanItem{Event: "err", Data: fmt.Sprintf("Failed to create request: %v", err)}
		return err // Return error early
	}

	// 2. Build HTTP Request
	httpReq, err := http.NewRequestWithContext(context.Background(), reqParams.Method, reqParams.URL, bytes.NewReader(reqParams.Body))
	if err != nil {
		common.BaseLogger.WithError(err).Error("Failed to create HTTP request")
		messageChan <- _type.MessageChanItem{Event: "err", Data: fmt.Sprintf("Failed to build request: %v", err)}
		return err
	}
	for k, v := range reqParams.Headers {
		httpReq.Header.Set(k, v)
	}

	// 3. Execute HTTP Request
	common.BaseLogger.Debugf("Sending stream request to: %s", reqParams.URL)
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		// Handle network/http error using provider
		handledErr := provider.HandleError(context.Background(), err, httpResp, nil)
		common.BaseLogger.WithError(handledErr).Error("Stream request failed (network/http)")
		messageChan <- _type.MessageChanItem{Event: "err", Data: handledErr.Error()}
		return handledErr
	}
	defer httpResp.Body.Close()

	// 4. Handle Non-OK HTTP Status Codes
	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(httpResp.Body)
		if readErr != nil {
			common.BaseLogger.WithError(readErr).Errorf("Failed to read error response body (status %d)", httpResp.StatusCode)
			bodyBytes = nil // Ensure bodyBytes is nil if read fails
		}
		// Use provider's HandleError
		handledErr := provider.HandleError(context.Background(), nil, httpResp, bodyBytes)
		common.BaseLogger.WithError(handledErr).Errorf("Stream request failed (status %d)", httpResp.StatusCode)
		messageChan <- _type.MessageChanItem{Event: "err", Data: handledErr.Error()}
		return handledErr
	}

	// 5. Process Stream Response
	common.BaseLogger.Debug("Stream request successful, processing response...")
	scanner := bufio.NewScanner(httpResp.Body)
	// Use standard SSE scanning
	scanner.Split(sse.MakeSplit(string(chain.Provider.Type)))
	isInitSent := false

	for scanner.Scan() {
		// Check context cancellation frequently within the loop
		select {
		case <-context.Background().Done():
			common.BaseLogger.Info("Context cancelled during stream processing.")
			messageChan <- _type.MessageChanItem{Event: "err", Data: context.Background().Err().Error()}
			return context.Background().Err()
		default:
			// Continue processing
		}

		eventBytes := scanner.Bytes()
		common.BaseLogger.Debugf("Received SSE event raw data: %s", string(eventBytes))

		// Skip empty lines or comment lines (though ScanSSE should handle comments)
		if len(bytes.TrimSpace(eventBytes)) == 0 {
			continue
		}

		// OPENROUTER specific keep-alive message
		if bytes.Equal(eventBytes, []byte(": OPENROUTER PROCESSING")) {
			common.BaseLogger.Debug("Skipping OpenRouter keep-alive message.")
			continue
		}

		// Parse the event data using the provider
		item, isDone, parseErr := provider.ParseStreamEvent(context.Background(), eventBytes)

		if parseErr != nil {
			common.BaseLogger.WithError(parseErr).Error("Failed to parse stream event")
			messageChan <- _type.MessageChanItem{Event: "err", Data: parseErr.Error()}
			return parseErr // Stop processing on parse error
		}

		if isDone {
			common.BaseLogger.Debug("Stream finished gracefully ([DONE] or equivalent detected).")
			break // Exit the loop normally
		}

		if item == nil {
			// Valid event, but no data item (e.g., empty event, comment handled by provider)
			continue
		}

		// Validate the parsed item (basic checks)
		if validationErr := item.IsValid(); validationErr != nil {
			// Check if it's an embedded error within the stream item itself
			if apiErr, ok := validationErr.(*_type.Error); ok {
				common.BaseLogger.WithError(apiErr).Warn("Stream item contains an API error")
				messageChan <- _type.MessageChanItem{Event: "err", Data: apiErr.Error()}
				return apiErr // Stop processing
			} else {
				// Other validation error (e.g., missing ID)
				common.BaseLogger.WithError(validationErr).Warn("Invalid stream item received")
				messageChan <- _type.MessageChanItem{Event: "err", Data: validationErr.Error()}
				return validationErr // Stop processing
			}
		}

		// Handle Azure prompt filter results (valid but skipped for data/init)
		if len(item.PromptFilterResults) > 0 && item.Id == "" {
			common.BaseLogger.Info("Received Azure prompt filter result, skipping data/init event.")
			continue
		}

		// Send init event for the first valid data item
		if !isInitSent {
			common.BaseLogger.Debugf("Sending init event: %+v", *item)
			messageChan <- _type.MessageChanItem{Event: "init", Data: *item}
			isInitSent = true
		}

		// Send data event if item has content
		if len(item.Choices) > 0 && item.Choices[0].Delta != nil && item.Choices[0].Delta.Content != "" {
			content := item.Choices[0].Delta.Content
			common.BaseLogger.Debugf("Sending data event: %s", content)
			messageChan <- _type.MessageChanItem{Event: "data", Data: content}
		} else {
			common.BaseLogger.Debugf("Stream item has no content to send: %+v", *item)
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		// Handle scanner error (e.g., connection closed unexpectedly)
		handledErr := provider.HandleError(context.Background(), scanErr, httpResp, nil)
		common.BaseLogger.WithError(handledErr).Error("Stream scanner error")
		messageChan <- _type.MessageChanItem{Event: "err", Data: handledErr.Error()}
		return handledErr
	}

	// Send finish event upon normal completion
	common.BaseLogger.Debug("Sending finish event.")
	messageChan <- _type.MessageChanItem{Event: "finish"}
	return nil // Indicate successful completion of stream processing
}
