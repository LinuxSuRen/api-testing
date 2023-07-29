package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-orm/pkg"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	opt := &option{}
	cmd := &cobra.Command{
		Use:   "store-orm",
		Short: "Storage extension of api-testing",
		RunE:  opt.runE,
	}
	flags := cmd.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7071, "The port of gRPC server")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var removeServer remote.LoaderServer
	if removeServer, err = NewRemoteServer(); err != nil {
		return
	}

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	err = gRPCServer.Serve(lis)
	return
}

type option struct {
	port int
}

type dbserver struct {
	remote.UnimplementedLoaderServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer() (s remote.LoaderServer, err error) {
	s = &dbserver{}
	return
}

func createDB(user, address, database string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:@tcp(%s)/%s?charset=utf8mb4", user, address, database)
	fmt.Println("try to connect to", dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		err = fmt.Errorf("failed to connect to %s, %v", dsn, err)
		return
	}

	db.AutoMigrate(&pkg.TestCase{})
	db.AutoMigrate(&pkg.TestSuite{})
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
		for key, val := range store.Properties {
			if key == "database" {
				database = val
			}
		}

		db, err = createDB(store.Username, store.URL, database)
		dbCache[store.Name] = db
	}
	return
}

func (s *dbserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	items := make([]*pkg.TestSuite, 0)

	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Find(&items)
	suites = &remote.TestSuites{}
	for i := range items {
		suites.Data = append(suites.Data, pkg.ConvertToGRPCTestSuite(items[i]))
	}
	return
}

func (s *dbserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Create(pkg.ConvertToDBTestSuite(testSuite))
	return
}

const nameQuery = `name = ?`

func (s *dbserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	query := &pkg.TestSuite{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	db.Find(&query, nameQuery, suite.Name)

	reply = pkg.ConvertToGRPCTestSuite(query)
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
	input := pkg.ConvertToDBTestSuite(suite)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}

	testSuiteIdentity(db, input).Updates(input)
	return
}

func testSuiteIdentity(db *gorm.DB, suite *pkg.TestSuite) *gorm.DB {
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
	items := make([]*pkg.TestCase, 0)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	db.Find(&items, "suite_name = ?", suite.Name)

	result = &server.TestCases{}
	for i := range items {
		result.Data = append(result.Data, pkg.ConvertToRemoteTestCase(items[i]))
	}
	return
}

func (s *dbserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	payload := pkg.ConverToDBTestCase(testcase)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	reply = &server.Empty{}
	db.Create(&payload)
	return
}

func (s *dbserver) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	item := &pkg.TestCase{}
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	db.Find(&item, "suite_name = ? AND name = ?", testcase.SuiteName, testcase.Name)

	result = pkg.ConvertToRemoteTestCase(item)
	return
}

func (s *dbserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	input := pkg.ConverToDBTestCase(testcase)
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
	input := pkg.ConverToDBTestCase(testcase)
	var db *gorm.DB
	if db, err = s.getClient(ctx); err != nil {
		return
	}
	testCaseIdentiy(db, input).Delete(input)
	return
}

func testCaseIdentiy(db *gorm.DB, testcase *pkg.TestCase) *gorm.DB {
	return db.Model(testcase).Where(fmt.Sprintf("suite_name = '%s' AND name = '%s'", testcase.SuiteName, testcase.Name))
}
