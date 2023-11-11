/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbserver struct {
	remote.UnimplementedLoaderServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer() (s remote.LoaderServer) {
	s = &dbserver{}
	return
}

func createDB(user, password, address, database, driver string) (db *gorm.DB, err error) {
	var dsn string
	switch driver {
	case "mysql", "":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", user, password, address, database)
	case "postgres":
		obj := strings.Split(address, ":")
		host, port := obj[0], obj[1]
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, database, port)
	case "clickhouse":
		dsn = fmt.Sprintf("tcp://%s?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20", address, database, user, password)
	default:
		err = fmt.Errorf("invalid database driver %q", driver)
		return
	}

	log.Printf("try to connect to %q", dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		err = fmt.Errorf("failed to connect to %s, %v", dsn, err)
		return
	}

	db.AutoMigrate(&TestCase{})
	db.AutoMigrate(&TestSuite{})
	return
}

var dbCache map[string]*gorm.DB = make(map[string]*gorm.DB)

func (s *dbserver) getClient(ctx context.Context) (db *gorm.DB, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to database")
	} else {
		var ok bool
		if db, ok = dbCache[store.Name]; ok && db != nil {
			return
		}

		database := "atest"
		driver := "mysql"
		if v, ok := store.Properties["database"]; ok && v != "" {
			database = v
		}
		if v, ok := store.Properties["driver"]; ok && v != "" {
			driver = v
		}

		if db, err = createDB(store.Username, store.Password, store.URL, database, driver); err == nil {
			dbCache[store.Name] = db
		}
	}
	return
}

func (s *dbserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	items := make([]*TestSuite, 0)

	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Find(&items)
	suites = &remote.TestSuites{}
	for i := range items {
		suites.Data = append(suites.Data, ConvertToGRPCTestSuite(items[i]))
	}
	return
}

func (s *dbserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Create(ConvertToDBTestSuite(testSuite))
	return
}

const nameQuery = `name = ?`

func (s *dbserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	query := &TestSuite{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Find(&query, nameQuery, suite.Name)

	reply = ConvertToGRPCTestSuite(query)
	if suite.Full {
		var testcases *server.TestCases
		if testcases, err = s.ListTestCases(ctx, &remote.TestSuite{
			Name: suite.Name,
		}); err == nil && testcases != nil {
			reply.Items = testcases.Data
		}
	}
	return
}

func (s *dbserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	input := ConvertToDBTestSuite(suite)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	testSuiteIdentity(db, input).Updates(input)
	return
}

func testSuiteIdentity(db *gorm.DB, suite *TestSuite) *gorm.DB {
	return db.Model(suite).Where(nameQuery, suite.Name)
}

func (s *dbserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Delete(suite, nameQuery, suite.Name)
	return
}

func (s *dbserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *server.TestCases, err error) {
	items := make([]*TestCase, 0)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	db.Find(&items, "suite_name = ?", suite.Name)

	result = &server.TestCases{}
	for i := range items {
		result.Data = append(result.Data, ConvertToRemoteTestCase(items[i]))
	}
	return
}

func (s *dbserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	payload := ConverToDBTestCase(testcase)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	reply = &server.Empty{}
	db.Create(&payload)
	return
}

func (s *dbserver) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	item := &TestCase{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	db.Find(&item, "suite_name = ? AND name = ?", testcase.SuiteName, testcase.Name)

	result = ConvertToRemoteTestCase(item)
	return
}

func (s *dbserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	input := ConverToDBTestCase(testcase)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	testCaseIdentiy(db, input).Updates(input)

	data := make(map[string]interface{})
	if input.ExpectBody == "" {
		data["expect_body"] = ""
	}
	if input.ExpectSchema == "" {
		data["expect_schema"] = ""
	}

	if len(data) > 0 {
		testCaseIdentiy(db, input).Updates(data)
	}
	return
}

func (s *dbserver) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	input := ConverToDBTestCase(testcase)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	testCaseIdentiy(db, input).Delete(input)
	return
}

func (s *dbserver) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	db, clientErr := s.getClient(ctx)
	reply = &server.ExtensionStatus{
		Ready:   err == nil && db != nil,
		Message: util.OKOrErrorMessage(clientErr),
		Version: version.GetVersion(),
	}
	return
}

func testCaseIdentiy(db *gorm.DB, testcase *TestCase) *gorm.DB {
	return db.Model(testcase).Where(fmt.Sprintf("suite_name = '%s' AND name = '%s'", testcase.SuiteName, testcase.Name))
}
