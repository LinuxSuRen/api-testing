package main

import (
	"context"
	"fmt"
	"net"

	"github.com/linuxsuren/api-testing/extensions/store-orm/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	removeServer := NewRemoteServer("root", "127.0.0.1:4000", "test")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7071))
	if err != nil {
		fmt.Println(err)
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	gRPCServer.Serve(lis)
}

type server struct {
	remote.UnimplementedLoaderServer
	db *gorm.DB
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(user, address, database string) remote.LoaderServer {
	db := createDB(user, address, database)
	return &server{db: db}
}

func createDB(user, address, database string) *gorm.DB {
	dsn := fmt.Sprintf("%s:@tcp(%s)/%s?charset=utf8mb4", user, address, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&pkg.TestCase{})
	db.AutoMigrate(&remote.TestSuite{})
	return db
}

func (s *server) ListTestSuite(context.Context, *remote.Empty) (suites *remote.TestSuites, err error) {
	suites = &remote.TestSuites{}
	items := make([]*remote.TestSuite, 23)
	s.db.Find(&items)
	suites.Data = items
	return
}

func (s *server) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	s.db.Create(testSuite)
	return
}

func (s *server) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	s.db.Find(&reply, "name = ?", suite.Name)
	return
}

func (s *server) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	testSuiteIdentity(s.db, suite).Updates(suite)
	return
}

func testSuiteIdentity(db *gorm.DB, suite *remote.TestSuite) *gorm.DB {
	return db.Model(&remote.TestSuite{}).Where("name = ?", suite.Name)
}

func (s *server) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	s.db.Delete(suite, "name = ?", suite.Name)
	return
}

func (s *server) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *remote.TestCases, err error) {
	items := make([]*pkg.TestCase, 0)
	s.db.Find(&items, "suite_name = ?", suite.Name)
	fmt.Println(items)

	result = &remote.TestCases{}
	for i := range items {
		result.Data = append(result.Data, &remote.TestCase{
			Name: items[i].Name,
		})
	}
	return
}

func (s *server) CreateTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.Empty, err error) {
	payload := pkg.ConverToDBTestCase(testcase)
	s.db.Create(&payload)
	return
}

func (s *server) GetTestCase(ctx context.Context, testcase *remote.TestCase) (result *remote.TestCase, err error) {
	fmt.Println(testcase.Name, testcase.SuiteName)

	item := &pkg.TestCase{}
	s.db.Find(&item, "suite_name = ? AND name = ?", testcase.SuiteName, testcase.Name)

	result = pkg.ConvertToRemoteTestCase(item)
	return
}

func (s *server) UpdateTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.TestCase, err error) {
	reply = &remote.TestCase{}
	fmt.Println(testcase.Request.Header)
	testCaseIdentiy(s.db, testcase).Updates(pkg.ConverToDBTestCase(testcase))
	return
}

func (s *server) DeleteTestCase(ctx context.Context, testcase *remote.TestCase) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	testCaseIdentiy(s.db, testcase).Delete(pkg.ConverToDBTestCase(testcase))
	return
}

func testCaseIdentiy(db *gorm.DB, testcase *remote.TestCase) *gorm.DB {
	return db.Model(&pkg.TestCase{}).Where(fmt.Sprintf("suite_name = '%s' AND name = '%s'", testcase.SuiteName, testcase.Name))
}
