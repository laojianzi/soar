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

package env

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/laojianzi/soar/common"
	"github.com/laojianzi/soar/database"

	"github.com/go-sql-driver/mysql"
	"github.com/kr/pretty"
)

var update = flag.Bool("update", false, "update .golden files")
var vEnv *VirtualEnv
var rEnv *database.Connector

func TestMain(m *testing.M) {
	// 初始化 init
	if common.DevPath == "" {
		_, file, _, _ := runtime.Caller(0)
		common.DevPath, _ = filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	}
	common.BaseDir = common.DevPath
	err := common.ParseConfig("")
	common.LogIfError(err, "init ParseConfig")
	common.Log.Debug("env_test init")
	vEnv, rEnv = BuildEnv()
	if _, err = vEnv.Version(); err != nil {
		fmt.Println(err.Error(), ", By pass all advisor test cases")
		os.Exit(0)
	}

	if _, err := rEnv.Version(); err != nil {
		fmt.Println(err.Error(), ", By pass all advisor test cases")
		os.Exit(0)
	}

	// 分割线
	flag.Parse()
	m.Run()

	// 环境清理
	vEnv.CleanUp()
}

func TestNewVirtualEnv(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	testSQL := []string{
		"use sakila",
		"select frm syntaxError",
		"CREATE TABLE t(id INT,c1 VARCHAR(20),PRIMARY KEY (id));",
		"ALTER TABLE t ADD INDEX `idx_c1`(c1);",
		"SELECT * FROM city WHERE country_id = 44;",
		"SELECT * FROM address WHERE address2 IS NOT NULL;",
		"SELECT * FROM address WHERE address2 IS NULL;",
		"SELECT * FROM address WHERE address2 >= 44;",
		"SELECT * FROM city WHERE country_id BETWEEN 44 AND 107;",
		"SELECT * FROM city WHERE city LIKE 'Ad%';",
		"SELECT * FROM city WHERE city = 'Aden' AND country_id = 107;",
		"SELECT * FROM city WHERE country_id > 31 AND city = 'Aden';",
		"SELECT * FROM address WHERE address_id > 8 AND city_id < 400 AND district = 'Nantou';",
		"SELECT * FROM address WHERE address_id > 8 AND city_id < 400;",
		"SELECT * FROM actor WHERE last_update='2006-02-15 04:34:33' AND last_name='CHASE' GROUP BY first_name;",
		"SELECT * FROM address WHERE last_update >='2014-09-25 22:33:47' GROUP BY district;",
		"SELECT * FROM address GROUP BY address,district;",
		"SELECT * FROM address WHERE last_update='2014-09-25 22:30:27' GROUP BY district,(address_id+city_id);",
		"SELECT * FROM customer WHERE active=1 ORDER BY last_name LIMIT 10;",
		"SELECT * FROM customer ORDER BY last_name LIMIT 10;",
		"SELECT * FROM customer WHERE address_id > 224 ORDER BY address_id LIMIT 10;",
		"SELECT * FROM customer WHERE address_id < 224 ORDER BY address_id LIMIT 10;",
		"SELECT * FROM customer WHERE active=1 ORDER BY last_name;",
		"SELECT * FROM customer WHERE address_id > 224 ORDER BY address_id;",
		"SELECT * FROM customer WHERE address_id IN (224,510) ORDER BY last_name;",
		"SELECT city FROM city WHERE country_id = 44;",
		"SELECT city,city_id FROM city WHERE country_id = 44 AND last_update='2006-02-15 04:45:25';",
		"SELECT city FROM city WHERE country_id > 44 AND last_update > '2006-02-15 04:45:25';",
		"SELECT * FROM city WHERE country_id=1 AND city='Kabul' ORDER BY last_update;",
		"SELECT * FROM city WHERE country_id>1 AND city='Kabul' ORDER BY last_update;",
		"SELECT * FROM city WHERE city_id>251 ORDER BY last_update; ",
		"SELECT * FROM city i INNER JOIN country o ON i.country_id=o.country_id;",
		"SELECT * FROM city i LEFT JOIN country o ON i.city_id=o.country_id;",
		"SELECT * FROM city i RIGHT JOIN country o ON i.city_id=o.country_id;",
		"SELECT * FROM city i LEFT JOIN country o ON i.city_id=o.country_id WHERE o.country_id IS NULL;",
		"SELECT * FROM city i RIGHT JOIN country o ON i.city_id=o.country_id WHERE i.city_id IS NULL;",
		"SELECT * FROM city i LEFT JOIN country o ON i.city_id=o.country_id UNION SELECT * FROM city i RIGHT JOIN country o ON i.city_id=o.country_id;",
		"SELECT * FROM city i LEFT JOIN country o ON i.city_id=o.country_id WHERE o.country_id IS NULL UNION SELECT * FROM city i RIGHT JOIN country o ON i.city_id=o.country_id WHERE i.city_id IS NULL;",
		"SELECT first_name,last_name,email FROM customer NATURAL LEFT JOIN address;",
		"SELECT first_name,last_name,email FROM customer NATURAL LEFT JOIN address;",
		"SELECT first_name,last_name,email FROM customer NATURAL RIGHT JOIN address;",
		"SELECT first_name,last_name,email FROM customer STRAIGHT_JOIN address ON customer.address_id=address.address_id;",
		"SELECT ID,name FROM (SELECT address FROM customer_list WHERE SID=1 ORDER BY phone LIMIT 50,10) a JOIN customer_list l ON (a.address=l.address) JOIN city c ON (c.city=l.city) ORDER BY phone DESC;",
	}

	err := common.GoldenDiff(func() {
		for _, sql := range testSQL {
			vEnv.BuildVirtualEnv(rEnv, sql)
			switch err := vEnv.Error.(type) {
			case nil:
				pretty.Println(sql, "OK")
			case *mysql.MySQLError:
				if err.Number != 1061 {
					t.Error(err)
				}
			case error:
				// unexpected EOF
				// 测试环境无法访问，或者被Disable的时候会进入这个分支
				pretty.Println(sql, err)
			default:
				t.Error(err)
			}
		}
	}, t.Name(), update)
	if err != nil {
		t.Error(err)
	}
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func TestCleanupTestDatabase(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	if common.Config.TestDSN.Disable {
		common.Log.Warn("common.Config.TestDSN.Disable=true, by pass TestCleanupTestDatabase")
		return
	}
	_, err := vEnv.Query("DROP DATABASE IF EXISTS optimizer_060102150405_xxxxxxxxxxxxxxxx")
	if err != nil {
		t.Error(err)
	}

	_, err = vEnv.Query("CREATE DATABASE optimizer_060102150405_xxxxxxxxxxxxxxxx")
	if err != nil {
		t.Error(err)
	}

	vEnv.CleanupTestDatabase()
	_, err = vEnv.Query("show create database optimizer_060102150405_xxxxxxxxxxxxxxxx")
	if err == nil {
		t.Error("optimizer_060102150405_xxxxxxxxxxxxxxxx exist, should be dropped")
	}

	_, err = vEnv.Query("DROP DATABASE IF EXISTS optimizer_060102150405")
	if err != nil {
		t.Error(err)
	}

	_, err = vEnv.Query("CREATE DATABASE optimizer_060102150405")
	if err != nil {
		t.Error(err)
	}

	vEnv.CleanupTestDatabase()
	_, err = vEnv.Query("DROP DATABASE optimizer_060102150405")
	if err != nil {
		t.Error("optimizer_060102150405 not exist, should not be dropped")
	}
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func TestGenTableColumns(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())

	pretty.Println(common.Config.TestDSN.Disable)
	if common.Config.TestDSN.Disable {
		common.Log.Warn("common.Config.TestDSN.Disable=true, by pass TestGenTableColumns")
		return
	}

	// 只能对sakila数据库进行测试
	if rEnv.Database == "sakila" {
		testSQL := []string{
			"SELECT * FROM city WHERE country_id = 44;",
			"SELECT country_id FROM city WHERE country_id = 44;",
			"SELECT country_id FROM city WHERE country_id > 44;",
		}

		metaList := []common.Meta{
			{
				"": &common.DB{
					Table: map[string]*common.Table{
						"city": common.NewTable("city"),
					},
				},
			},
			{
				"sakila": &common.DB{
					Table: map[string]*common.Table{
						"city": common.NewTable("city"),
					},
				},
			},
			{
				"sakila": &common.DB{
					Table: map[string]*common.Table{
						"city": {
							TableName: "city",
							Column: map[string]*common.Column{
								"country_id": {
									Name: "country_id",
								},
							},
						},
					},
				},
			},
		}

		for i, sql := range testSQL {
			vEnv.BuildVirtualEnv(rEnv, sql)
			tFlag := false
			columns := vEnv.GenTableColumns(metaList[i])
			if _, ok := columns["sakila"]; ok {
				if _, okk := columns["sakila"]["city"]; okk {
					if length := len(columns["sakila"]["city"]); length >= 1 {
						tFlag = true
					}
				}
			}

			if !tFlag {
				t.Errorf("columns: \n%s", pretty.Sprint(columns))
			}
		}
	}
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func TestCreateTable(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	orgSamplingCondition := common.Config.SamplingCondition
	common.Config.SamplingCondition = "LIMIT 1"

	orgREnvDatabase := rEnv.Database
	rEnv.Database = "sakila"
	tables := []string{
		"actor",
		"address",
		"category",
		"city",
		"country",
		"customer",
		"film",
		"film_actor",
		"film_category",
		"film_text",
		"inventory",
		"language",
		"payment",
		"rental",
		"staff",
		"store",
		"staff_list",
		"customer_list",
		"actor_info",
		"sales_by_film_category",
		"sales_by_store",
		"nicer_but_slower_film_list",
		"film_list",
	}
	for _, table := range tables {
		err := vEnv.createTable(rEnv, table)
		if err != nil {
			t.Error(err)
		}
	}
	common.Config.SamplingCondition = orgSamplingCondition
	rEnv.Database = orgREnvDatabase
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}

func TestCreateDatabase(t *testing.T) {
	common.Log.Debug("Entering function: %s", common.GetFunctionName())
	orgREnvDatabase := rEnv.Database
	rEnv.Database = "sakila"
	err := vEnv.createDatabase(rEnv)
	if err != nil {
		t.Error(err)
	}
	if vEnv.DBHash("sakila") == "sakila" {
		t.Errorf("database: sakila rehashed failed!")
	}

	if vEnv.DBHash("not_exist_db") != "not_exist_db" {
		t.Errorf("database: not_exist_db rehashed!")
	}
	rEnv.Database = orgREnvDatabase
	common.Log.Debug("Exiting function: %s", common.GetFunctionName())
}
