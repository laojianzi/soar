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

package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/tidwall/gjson"
)

// GoldenDiff 从 gofmt 学来的测试方法
// https://medium.com/soon-london/testing-with-golden-files-in-go-7fccc71c43d3
func GoldenDiff(f func(), name string, update *bool) error {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	str := captureOutput(f)
	_, err := w.WriteString(str)
	if err != nil {
		Log.Warning(err.Error())
	}
	err = w.Flush()
	if err != nil {
		Log.Warning(err.Error())
	}

	gp := filepath.Join("testdata", name+".golden")
	got := bytes.ReplaceAll(b.Bytes(), []byte{'\r', '\n'}, []byte{'\n'})
	if *update {
		if err = ioutil.WriteFile(gp, got, 0644); err != nil {
			err = fmt.Errorf("%s failed to update golden file: %s", name, err)
			return err
		}
	}

	want, err := ioutil.ReadFile(gp)
	if err != nil {
		err = fmt.Errorf("%s failed reading .golden: %s", name, err)
	}

	want = bytes.ReplaceAll(want, []byte{'\r', '\n'}, []byte{'\n'})
	if !bytes.Equal(got, want) {
		err = fmt.Errorf("%s does not match .golden file\nwant: \n%q\ngot: \n%q", name, string(want), string(got))
	}

	return err
}

// captureOutput 获取函数标准输出
func captureOutput(f func()) string {
	// keep backup of the real stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// copy the output in a separate goroutine so printing can't block indefinitely
	outC := make(chan string)
	go func() {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			panic(err)
		}
		outC <- string(buf)
	}()

	// execute function
	f()

	// back to normal state
	if err := w.Close(); err != nil {
		panic(err)
	}
	os.Stdout = oldStdout // restoring the real stdout
	out := <-outC
	os.Stdout = oldStdout
	return out
}

// SortedKey sort map[string]interface{}, use in range clause
func SortedKey(m interface{}) []string {
	var keys []string
	switch reflect.TypeOf(m).Kind() {
	case reflect.Map:
		switch reflect.TypeOf(m).Key().Kind() {
		case reflect.String:
			for _, k := range reflect.ValueOf(m).MapKeys() {
				keys = append(keys, k.String())
			}
		}
	}
	sort.Strings(keys)
	return keys
}

// jsonFind internal function
func jsonFind(json string, name string, find *[]string) (next []string) {
	res := gjson.Parse(json)
	res.ForEach(func(key, value gjson.Result) bool {
		if key.String() == name {
			*find = append(*find, value.String())
		}
		switch value.Type {
		case gjson.Number, gjson.True, gjson.False, gjson.Null:
		default:
			// String, JSON
			next = append(next, value.String())
		}
		return true // keep iterating
	})
	return next
}

// JSONFind iterate find name in json
// TODO: for complicate SQL JSONFind will run a long time for json interactions.
func JSONFind(json string, name string) []string {
	Log.Debug("Entering function: %s", GetFunctionName())
	var find []string
	next := []string{json}
	for len(next) > 0 {
		var tmpNext []string
		for _, subJSON := range next {
			tmpNext = append(tmpNext, jsonFind(subJSON, name, &find)...)
		}
		next = tmpNext
	}
	Log.Debug("Exiting function: %s", GetFunctionName())
	return find
}

// RemoveDuplicatesItem remove duplicate item from list
func RemoveDuplicatesItem(duplicate []string) []string {
	m := make(map[string]bool)
	for _, item := range duplicate {
		if _, ok := m[item]; !ok {
			m[item] = true
		}
	}

	var unique []string
	for item := range m {
		unique = append(unique, item)
	}
	sort.Strings(unique)
	return unique
}
