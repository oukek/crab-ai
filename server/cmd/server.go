package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"oukek/crab-ai/app/router"
)

var graceful *bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务",
	Long: `启动服务
`,
	Run: func(cmd *cobra.Command, args []string) {
		savePid()
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	graceful = serverCmd.Flags().Bool("graceful", false, "灰度重启")
}

var listener net.Listener
var server *http.Server

func startServer() {

	var err error

	r := router.GetRouter()
	server = &http.Server{
		Addr:    os.Getenv("LISTEN_ADDRESS"),
		Handler: r,
	}

	if *graceful {
		log.Print("main: Listening to existing file descriptor 3.")
		// cmd.ExtraFiles: If non-nil, entry i becomes file descriptor 3+i.
		// when we put socket FD at the first entry, it will always be 3(0+3)
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		log.Print("main: Listening on a new file descriptor.")
		listener, err = net.Listen("tcp", server.Addr)
	}

	if err != nil {
		log.Fatalf("listener error: %v", err)
	}

	go func() {
		// 服务连接
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	signalHandler()

}

// 将当前监听的socket传给子进程
func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}

	f, err := tl.File()
	if err != nil {
		return err
	}

	args := []string{"server", "--graceful"}
	cmd := exec.Command(os.Args[0], args...)
	fmt.Printf("%v\n", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}

func savePid() {
	pidfile := "./novel_server.pid"
	pf, err := os.OpenFile(pidfile, os.O_WRONLY|os.O_TRUNC, 0)
	defer pf.Close()
	if os.IsNotExist(err) {
		pf, err = os.Create(pidfile)
		if err != nil {
			log.Fatal("create pid file error.")
			return
		}
	}
	pid := os.Getpid()
	_, err = pf.Write([]byte(fmt.Sprintf("%d", pid)))
	if err != nil {
		log.Fatal("write pid failed.")
	}
}
