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
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3API interface {
	ListObjectsWithContext(ctx aws.Context, input *s3.ListObjectsInput, opts ...request.Option) (*s3.ListObjectsOutput, error)
	PutObjectWithContext(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
	GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
	DeleteObjectWithContext(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error)
}

type S3Creator interface {
	New(p client.ConfigProvider, cfgs ...*aws.Config) S3API
}

type DefaultS3Creator struct{}

func (d *DefaultS3Creator) New(p client.ConfigProvider, cfgs ...*aws.Config) S3API {
	return s3.New(p, cfgs...)
}

type fakeS3 struct {
	data map[*string][]byte
}

func (f *fakeS3) New(p client.ConfigProvider, cfgs ...*aws.Config) S3API {
	return f
}

func (f *fakeS3) ListObjectsWithContext(ctx aws.Context, input *s3.ListObjectsInput, opts ...request.Option) (output *s3.ListObjectsOutput, err error) {
	output = &s3.ListObjectsOutput{}
	for k := range f.data {
		output.Contents = append(output.Contents, &s3.Object{
			Key: k,
		})
	}
	return
}
func (f *fakeS3) PutObjectWithContext(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
	data, err := io.ReadAll(input.Body)
	f.data[input.Key] = data
	return nil, err
}
func (f *fakeS3) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (output *s3.GetObjectOutput, err error) {
	for k := range f.data {
		if *input.Key == *k {
			output = &s3.GetObjectOutput{
				Body: io.NopCloser(bytes.NewReader(f.data[k])),
			}
			break
		}
	}
	return
}
func (f *fakeS3) DeleteObjectWithContext(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
	for k := range f.data {
		if *input.Key == *k {
			delete(f.data, k)
			break
		}
	}
	return nil, nil
}
