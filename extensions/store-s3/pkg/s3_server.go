package pkg

import (
	"bytes"
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"gopkg.in/yaml.v3"
)

type s3Client struct {
	remote.UnimplementedLoaderServer
}

func NewRemoteServer() (remote.LoaderServer, error) {
	return &s3Client{}, nil
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
	}); err == nil {
		var suite *testing.TestSuite
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
	}
	return
}
func (s *s3Client) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	suite := remote.ConvertToNormalTestSuite(testSuite)
	reply = &server.Empty{}

	var data []byte
	if data, err = yaml.Marshal(suite); err != nil {
		return
	}

	var client *s3WithBucket
	if client, err = s.getClient(ctx); err != nil {
		return
	}

	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(client.bucket),
		Key:    aws.String(suite.Name + ".yaml"),
		Body:   bytes.NewReader(data),
	})
	return
}
func (s *s3Client) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var client *s3WithBucket
	if client, err = s.getClient(ctx); err != nil || client == nil {
		return
	}

	var objOutput *s3.GetObjectOutput
	if objOutput, err = client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(client.bucket),
		Key:    aws.String(suite.Name + ".yaml"),
	}); err == nil {
		data := objOutput.Body

		var suite *testing.TestSuite
		if suite, err = testing.ParseFromStream(data); err == nil {
			reply = remote.ConvertToGRPCTestSuite(suite)
		}
	}
	return
}
func (s *s3Client) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var oldSuite *remote.TestSuite
	if oldSuite, err = s.GetTestSuite(ctx, suite); err != nil {
		return
	}

	suite.Items = oldSuite.Items
	_, err = s.CreateTestSuite(ctx, suite)
	return
}
func (s *s3Client) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var client *s3WithBucket
	if client, err = s.getClient(ctx); err != nil || client == nil {
		return
	}

	_, err = client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(client.bucket),
		Key:    aws.String(suite.Name + ".yaml"),
	})
	return
}
func (s *s3Client) ListTestCases(ctx context.Context, suite *remote.TestSuite) (result *server.TestCases, err error) {
	if suite, err = s.GetTestSuite(ctx, suite); err != nil {
		return
	}

	result = &server.TestCases{
		Data: suite.Items,
	}
	return
}
func (s *s3Client) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}

	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err != nil {
		return
	}

	suite.Items = append(suite.Items, testcase)
	_, err = s.CreateTestSuite(ctx, suite)
	return
}
func (s *s3Client) GetTestCase(ctx context.Context, testcase *server.TestCase) (result *server.TestCase, err error) {
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err != nil {
		return
	}

	for _, item := range suite.Items {
		if item.Name == testcase.Name {
			result = item
			break
		}
	}
	return
}
func (s *s3Client) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err != nil {
		return
	}

	for i, item := range suite.Items {
		if item.Name == testcase.Name {
			suite.Items[i] = testcase
			break
		}
	}

	_, err = s.CreateTestSuite(ctx, suite)
	return
}
func (s *s3Client) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var suite *remote.TestSuite
	if suite, err = s.GetTestSuite(ctx, &remote.TestSuite{
		Name: testcase.SuiteName,
	}); err != nil {
		return
	}

	for i, item := range suite.Items {
		if item.Name == testcase.Name {
			suite.Items = append(suite.Items[:i], suite.Items[i+1:]...)
			break
		}
	}

	_, err = s.UpdateTestSuite(ctx, suite)
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
		cred := credentials.NewStaticCredentials(options.AccessKeyID, options.SecretAccessKey, options.SessionToken)

		config := aws.Config{
			Region:           aws.String(options.Region),
			Endpoint:         aws.String(store.URL),
			DisableSSL:       aws.Bool(options.DisableSSL),
			S3ForcePathStyle: aws.Bool(options.ForcePathStyle),
			Credentials:      cred,
		}

		var sess *session.Session
		sess, err = session.NewSession(&config)
		if err != nil {
			return
		}

		svc := s3.New(sess)
		db = &s3WithBucket{S3: svc, bucket: options.Bucket}
		clientCache[store.Name] = db
	}
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
	*s3.S3
	bucket string
}

var clientCache map[string]*s3WithBucket = make(map[string]*s3WithBucket)
