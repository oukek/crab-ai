package cmd

import (
	"log"

	_ "embed"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"github.com/ncruces/go-sqlite3"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "测试功能",
	Long:  `测试各种功能`,
	Run: func(cmd *cobra.Command, args []string) {
		testVecVersion()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

// 测试vec_version函数
func testVecVersion() {
	// ctx := context.Background()

	// // 创建配置时启用原子特性
	// config := wazero.NewRuntimeConfig()
	// config.WithCoreFeatures(wasm.CoreFeaturesV2)
	// r := wazero.NewRuntimeWithConfig(ctx, config)
	// defer r.Close(ctx)
	// 创建内存数据库连接
	db, err := sqlite3.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 准备查询语句
	stmt, _, err := db.Prepare(`SELECT vec_version()`)
	if err != nil {
		log.Fatal(err)
	}

	// 执行查询
	stmt.Step()
	log.Printf("vec_version=%s\n", stmt.ColumnText(0))
	stmt.Close()
}
