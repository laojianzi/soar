package soar

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/XiaoMi/soar/common"
)

var update = flag.Bool("update", false, "update .golden files")

func TestMain(m *testing.M) {
	// 初始化 init
	if common.DevPath == "" {
		_, file, _, _ := runtime.Caller(0)
		common.DevPath, _ = filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	}
	common.BaseDir = common.DevPath
	err := common.ParseConfig("")
	common.LogIfError(err, "init ParseConfig")
	common.Log.Debug("mysql_test init")
	_ = update // check if var success init

	// 分割线
	flag.Parse()
	m.Run()

	// 环境清理
	//
}

func Test_Soar_initQuery(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	// direct query
	query := initQuery("select 1")
	if query != "select 1" {
		t.Errorf("want 'select 1', got %s", query)
	}

	// read from file
	initQuery(common.DevPath + "/README.md")

	orgStdin := os.Stdin
	tmpStdin, err := os.Open(common.DevPath + "/VERSION")
	if err != nil {
		t.Error(err)
	}
	os.Stdin = tmpStdin
	fmt.Println(initQuery(""))
	os.Stdin = orgStdin
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func Test_Soar_reportTool(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	orgRerportType := common.Config.ReportType
	types := []string{"html", "md2html", "explain-digest", "chardet", "remove-comment"}
	for _, tp := range types {
		common.Config.ReportType = tp
		fmt.Println(reportTool(tp, []byte{}))
	}
	common.Config.ReportType = orgRerportType
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func Test_Soar_helpTools(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())

	orgConfig := common.CheckConfig
	common.CheckConfig = true
	helpTools()
	common.CheckConfig = orgConfig

	orgConfig = common.PrintVersion
	common.PrintVersion = true
	helpTools()
	common.PrintVersion = orgConfig

	orgConfig = common.PrintConfig
	common.PrintConfig = true
	helpTools()
	common.PrintConfig = orgConfig

	orgConfig = common.Config.ListHeuristicRules
	common.Config.ListHeuristicRules = true
	helpTools()
	common.Config.ListHeuristicRules = orgConfig

	orgConfig = common.Config.ListRewriteRules
	common.Config.ListRewriteRules = true
	helpTools()
	common.Config.ListRewriteRules = orgConfig

	orgConfig = common.Config.ListTestSqls
	common.Config.ListTestSqls = true
	helpTools()
	common.Config.ListTestSqls = orgConfig

	orgConfig = common.Config.ListReportTypes
	common.Config.ListReportTypes = true
	helpTools()
	common.Config.ListReportTypes = orgConfig
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func Test_Soar_verboseInfo(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	orgVerbose := common.Config.Verbose
	common.Config.Verbose = true
	err := common.GoldenDiff(func() {
		// Syntax check OK
		orgSyntaxCheck := common.Config.OnlySyntaxCheck
		common.Config.OnlySyntaxCheck = true
		verboseInfo()
		common.Config.OnlySyntaxCheck = orgSyntaxCheck

		// MySQL environment verbose info
		orgTestDSNDisable := common.Config.TestDSN.Disable
		common.Config.TestDSN.Disable = true
		verboseInfo()
		common.Config.TestDSN.Disable = orgTestDSNDisable

		orgOnlineDSNDisable := common.Config.OnlineDSN.Disable
		common.Config.OnlineDSN.Disable = true
		verboseInfo()
		common.Config.OnlineDSN.Disable = orgOnlineDSNDisable
	}, t.Name(), update)
	if err != nil {
		t.Error(err)
	}

	common.Config.Verbose = orgVerbose
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}
