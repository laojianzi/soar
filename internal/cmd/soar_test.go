/*
 * Copyright 2018 Xiaomi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/laojianzi/soar/common"
)

var update = flag.Bool("update", false, "update .golden files")

func TestMain(m *testing.M) {
	// 初始化 init
	if common.DevPath == "" {
		_, file, _, _ := runtime.Caller(0)
		common.DevPath, _ = filepath.Abs(filepath.Dir(filepath.Join(file, "..", "..")))
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

func Test_Cmd_initQuery(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	// direct query
	query := initQuery("select 1")
	if query != "select 1" {
		t.Errorf("want 'select 1', got %s", query)
	}

	oldDevPath := common.DevPath
	if dirList := strings.Split(common.DevPath, string(os.PathSeparator)); dirList[len(dirList)-1] == "cmd" {
		common.DevPath = filepath.Join(common.DevPath, "..")
	}

	// read from file
	initQuery(filepath.Join(common.DevPath, "README.md"))

	orgStdin := os.Stdin
	tmpStdin, err := os.Open(filepath.Join(common.DevPath, "VERSION"))
	if err != nil {
		t.Error(err)
	}
	common.DevPath = oldDevPath

	os.Stdin = tmpStdin
	fmt.Println(initQuery(""))
	os.Stdin = orgStdin
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func Test_Cmd_reportTool(_ *testing.T) {
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

func Test_Cmd_helpTools(_ *testing.T) {
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

func Test_Cmd_verboseInfo(t *testing.T) {
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
