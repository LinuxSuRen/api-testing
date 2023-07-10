package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-orm/pkg"
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
	flags.StringVarP(&opt.user, "user", "u", "root", "The user name of database")
	flags.StringVarP(&opt.address, "address", "", "127.0.0.1:4000", "The address of database")
	flags.StringVarP(&opt.database, "database", "", "test", "The database name")
	flags.IntVarP(&opt.port, "port", "p", 7071, "The port of gRPC server")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var removeServer remote.LoaderServer
	if removeServer, err = NewRemoteServer(o.user, o.address, o.database); err != nil {
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
	user, address, database string
	port                    int
}

type server struct {
	remote.UnimplementedLoaderServer
	db *gorm.DB
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(user, address, database string) (s remote.LoaderServer, err error) {
	var db *gorm.DB
	if db, err = createDB(user, address, database); err == nil {
		s = &server{db: db}
	}
	return
}

func createDB(user, address, database string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:@tcp(%s)/%s?charset=utf8mb4", user, address, database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return
	}

	db.AutoMigrate(&pkg.TestCase{})
	db.AutoMigrate(&pkg.TestSuite{})
	return
}

func (s *server) ListTestSuite(context.Context, *remote.Empty) (suites *remote.TestSuites, err error) {
	items := make([]*pkg.TestSuite, 0)
	s.db.Find(&items)

	suites = &remote.TestSuites{}
	for i := range items {
		suites.Data = append(suites.Data, pkg.ConvertToGRPCTestSuite(items[i]))
	}
	return
}

func (s *server) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	s.db.Create(pkg.ConvertToDBTestSuite(testSuite))
	return
}

const nameQuery = `name = ?`

func (s *server) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	query := &pkg.TestSuite{}
	s.db.Find(&query, nameQuery, suite.Name)

	reply = pkg.ConvertToGRPCTestSuite(query)
	if suite.Full {
		var testcases *remote.TestCases
		if testcases, err = s.ListTestCases(ctx, &remote.TestSuite{
			Name: suite.Name,
		}); err == nil && testcases != nil {
			reply.Items = testcases.Data
		}
	}
	return
}

func (s *server) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	input := pkg.ConvertToDBTestSuite(suite)
	testSuiteIdentity(s.db, input).Updates(input)
	return
}

func testSuiteIdentity(db *gorm.DB, suite *pkg.TestSuite) *gorm.DB {
	return db.Model(suite).Where(nameQuery, suite.Name)
}

func (s *server) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	s.db.Delete(suite, nameQuery, suite.Name)
	return
}

func (s *server) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *remote.TestCases, err error) {
	items := make([]*pkg.TestCase, 0)
	s.db.Find(&items, "suite_name = ?", suite.Name)

	result = &remote.TestCases{}
	for i := range items {
		result.Data = append(result.Data, pkg.ConvertToRemoteTestCase(items[i]))
	}
	return
}

func (s *server) CreateTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.Empty, err error) {
	payload := pkg.ConverToDBTestCase(testcase)
	s.db.Create(&payload)
	return
}

func (s *server) GetTestCase(ctx context.Context, testcase *remote.TestCase) (result *remote.TestCase, err error) {
	item := &pkg.TestCase{}
	s.db.Find(&item, "suite_name = ? AND name = ?", testcase.SuiteName, testcase.Name)

	result = pkg.ConvertToRemoteTestCase(item)
	return
}

func (s *server) UpdateTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.TestCase, err error) {
	reply = &remote.TestCase{}
	input := pkg.ConverToDBTestCase(testcase)
	testCaseIdentiy(s.db, input).Updates(input)

	data := make(map[string]interface{})
	if input.ExpectBody == "" {
		data["expect_body"] = ""
	}
	if input.ExpectSchema == "" {
		data["expect_schema"] = ""
	}

	if len(data) > 0 {
		testCaseIdentiy(s.db, input).Updates(data)
	}
	return
}

func (s *server) DeleteTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	input := pkg.ConverToDBTestCase(testcase)
	testCaseIdentiy(s.db, input).Delete(input)
	return
}

func testCaseIdentiy(db *gorm.DB, testcase *pkg.TestCase) *gorm.DB {
	return db.Model(testcase).Where(fmt.Sprintf("suite_name = '%s' AND name = '%s'", testcase.SuiteName, testcase.Name))
}
