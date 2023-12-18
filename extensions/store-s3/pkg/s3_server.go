/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pkg

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
	"gopkg.in/yaml.v3"
)

type s3Client struct {
	S3Creator S3Creator
	remote.UnimplementedLoaderServer
}

func NewRemoteServer(s3Creator S3Creator) remote.LoaderServer {
	return &s3Client{S3Creator: s3Creator}
}

func (s *s3Client) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	suites = &remote.TestSuites{}
	var client *s3WithBucket
	if client, err = s.getClient(ctx); err != nil || client == nil {
		return
	}

	var list *s3.ListObjectsOutput
	if list, err = client.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(client.bucket),
	}); err == nil && list != nil {
		suites, err = listObjectsOutputToTestSuite(ctx, list, client)
	}
	return
}
func listObjectsOutputToTestSuite(ctx context.Context, list *s3.ListObjectsOutput, client *s3WithBucket) (
	suites *remote.TestSuites, err error) {
	var suite *testing.TestSuite
	suites = &remote.TestSuites{}
	for _, obj := range list.Contents {
		if !strings.HasSuffix(*obj.Key, ".yaml") {
			continue
		}

		var objOutput *s3.GetObjectOutput
		if objOutput, err = client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(client.bucket),
			Key:    obj.Key,
		}); err == nil {
			data := objOutput.Body
			if suite, err = testing.ParseFromStream(data); err == nil {
				suites.Data = append(suites.Data, remote.ConvertToGRPCTestSuite(suite))
			}
		}
	}
	return
}
func (s *s3Client) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	suite := remote.ConvertToNormalTestSuite(testSuite)
	reply = &server.Empty{}

	var data []byte
	if data, err = yaml.Marshal(suite); err == nil {
		var client *s3WithBucket
		if client, err = s.getClient(ctx); err == nil {
			_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
				Bucket: aws.String(client.bucket),
				Key:    generateKey(suite.Name),
				Body:   bytes.NewReader(data),
			})
		}
	}

	return
}
func (s *s3Client) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var client *s3WithBucket
	if client, err = s.getClient(ctx); err == nil && client != nil {
		var objOutput *s3.GetObjectOutput
		if objOutput, err = client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(client.bucket),
			Key:    generateKey(suite.Name),
		}); err == nil && objOutput != nil {
			data := objOutput.Body

			var suite *testing.TestSuite
			if suite, err = testing.ParseFromStream(data); err == nil {
				reply = remote.ConvertToGRPCTestSuite(suite)
			}
		}
	}
	return
}
func (s *s3Client) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var oldSuite *remote.TestSuite
	if oldSuite, err = s.GetTestSuite(ctx, suite); err == nil {
		suite.Items = oldSuite.Items
		_, err = s.CreateTestSuite(ctx, suite)
	}
	return
}
func (s *s3Client) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var client *s3WithBucket
	if client, err = s.getClient(ctx); err == nil && client != nil {
		_, err = client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(client.bucket),
			Key:    generateKey(suite.Name),
		})
	}
	return
}
func (s *s3Client) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *server.TestCases, err error) {
	if suite, err = s.GetTestSuite(ctx, suite); err == nil {
		result = &server.TestCases{
			Data: suite.Items,
		}
	}
	return
}
func (s *s3Client) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}

	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err == nil {
		suite.Items = append(suite.Items, testcase)
		_, err = s.CreateTestSuite(ctx, suite)
	}
	return
}
func (s *s3Client) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err == nil {
		result = getTestCaseByName(suite, testcase.Name)
	}
	return
}
func (s *s3Client) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err == nil {
		suite = updateTestCase(suite, testcase)
		_, err = s.CreateTestSuite(ctx, suite)
	}
	return
}
func (s *s3Client) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err == nil {
		suite = removeTestCaseByName(suite, testcase.Name)
		_, err = s.UpdateTestSuite(ctx, suite)
	}
	return
}
func (s *s3Client) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	_, clientErr := s.ListTestSuite(ctx, in)
	reply = &server.ExtensionStatus{
		Ready:   err == nil,
		Message: util.OKOrErrorMessage(clientErr),
		Version: version.GetVersion(),
	}
	return
}
func (s *s3Client) PProf(ctx context.Context, in *server.PProfRequest) (data *server.PProfData, err error) {
	log.Println("pprof", in.Name)

	data = &server.PProfData{
		Data: extension.LoadPProf(in.Name),
	}
	return
}
func (s *s3Client) getClient(ctx context.Context) (db *s3WithBucket, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to s3 server")
	} else {
		var ok bool
		if db, ok = clientCache[store.Name]; ok && db != nil {
			return
		}

		options := mapToS3Options(store.Properties)
		if options.AccessKeyID == "" {
			options.AccessKeyID = store.Username
		}
		if options.SecretAccessKey == "" {
			options.SecretAccessKey = store.Password
		}

		log.Println("s3 server", store.URL)

		var sess *session.Session
		sess, err = createClientFromSs3Options(options, store.URL)
		if err == nil {
			svc := s.S3Creator.New(sess)
			db = &s3WithBucket{S3API: svc, bucket: options.Bucket}
			clientCache[store.Name] = db
		}
	}
	return
}
func createClientFromSs3Options(options s3Options, storeURL string) (sess *session.Session, err error) {
	cred := credentials.NewStaticCredentials(options.AccessKeyID, options.SecretAccessKey, options.SessionToken)

	config := aws.Config{
		Region:           aws.String(options.Region),
		Endpoint:         aws.String(storeURL),
		DisableSSL:       aws.Bool(options.DisableSSL),
		S3ForcePathStyle: aws.Bool(options.ForcePathStyle),
		Credentials:      cred,
	}

	sess, err = session.NewSession(&config)
	return
}

func mapToS3Options(data map[string]string) (opt s3Options) {
	opt.AccessKeyID = data["accesskeyid"]
	opt.SecretAccessKey = data["secretaccesskey"]
	opt.SessionToken = data["sessiontoken"]
	opt.Region = data["region"]
	opt.DisableSSL = data["disablessl"] == "true"
	opt.ForcePathStyle = data["forcepathstyle"] == "true"
	opt.Bucket = data["bucket"]
	return
}

func generateKey(name string) *string {
	return aws.String(name + ".yaml")
}

func removeTestCaseByName(suite *remote.TestSuite, name string) *remote.TestSuite {
	for i, item := range suite.Items {
		if item.Name == name {
			suite.Items = append(suite.Items[:i], suite.Items[i+1:]...)
			break
		}
	}
	return suite
}

func updateTestCase(suite *remote.TestSuite, testcase *server.TestCase) *remote.TestSuite {
	for i, item := range suite.Items {
		if item.Name == testcase.Name {
			suite.Items[i] = testcase
			break
		}
	}
	return suite
}

func getTestCaseByName(suite *remote.TestSuite, name string) (result *server.TestCase) {
	for _, item := range suite.Items {
		if item.Name == name {
			result = item
			break
		}
	}
	return
}

type s3Options struct {
	// AWS Access key ID
	AccessKeyID string `yaml:"accessKeyID"`
	// AWS Secret Access Key
	SecretAccessKey string `yaml:"secretAccessKey"`
	// AWS Session Token
	SessionToken string `yaml:"sessionToken"`
	// AWS Region
	Region         string `yaml:"region"`
	DisableSSL     bool   `yaml:"disableSSL"`
	ForcePathStyle bool   `yaml:"forcePathStyle"`
	Bucket         string `yaml:"bucket"`
}

type s3WithBucket struct {
	S3API
	bucket string
}

var clientCache map[string]*s3WithBucket = make(map[string]*s3WithBucket)
