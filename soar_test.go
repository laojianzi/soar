package soar

import (
	"strconv"
	"strings"
	"testing"

	"github.com/XiaoMi/soar/common"
)

func TestRun(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	// rules, out := Run("SELECT * FROM film") // 返回 L1 和 L4
	rules, out := Run("SELECT * FROM film WHERE title = 'abc'") // 返回 L1
	t.Log(out)
	// 如大于 L2 的无法接受
	for _, rule := range rules {
		severityLevel, err := strconv.Atoi(strings.Trim(rule.Severity, "L"))
		if err != nil {
			t.Fatal(err)
		}

		if severityLevel > 2 {
			t.Errorf("%+v", rule)
		}
	}

	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}
