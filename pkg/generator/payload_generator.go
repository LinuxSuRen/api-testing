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
package generator

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type grpcPayloadGenerator struct {
}

func NewGrpcPayloadGenerator() CodeGenerator {
	return &grpcPayloadGenerator{}
}

func (g *grpcPayloadGenerator) Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (result string, err error) {
	result, err = generateGRPCPayloadAsJSON(testSuite.Spec.RPC,
		parseGRPCService(strings.TrimPrefix(testcase.Request.API, testSuite.API)))
	return
}

func parseGRPCService(service string) string {
	service = strings.TrimPrefix(service, "/")
	return strings.ReplaceAll(service, "/", ".")
}

func generateGRPCPayloadAsJSON(rpc *testing.RPCDesc, service string) (resultJSON string, err error) {
	if rpc == nil {
		err = errors.New("rpc is missing")
		return
	}

	protoFile := rpc.ProtoFile
	protoContent := rpc.Raw

	if protoFile == "" {
		// don't really need a regular file, just give it a name
		protoFile = "placeholder.proto"
	}

	var (
		importPath     []string
		parentProtoDir string
	)
	protoFile, importPath, parentProtoDir, err = util.LoadProtoFiles(protoFile)
	if err != nil {
		return
	}

	if len(importPath) == 0 {
		importPath = rpc.ImportPath
	}

	if parentProtoDir != "" {
		for i, p := range importPath {
			importPath[i] = filepath.Join(parentProtoDir, p)
		}
		if len(importPath) == 0 {
			importPath = append(importPath, parentProtoDir)
		}
	}

	var protoLibrary map[string]string
	if protoLibrary, err = apispec.GetProtoFiles(); err != nil {
		return
	}

	if protoContent != "" {
		protoLibrary[protoFile] = protoContent
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(
			&protocompile.SourceResolver{
				Accessor: func(path string) (io.ReadCloser, error) {
					if content, ok := protoLibrary[strings.TrimPrefix(path, parentProtoDir+"/")]; ok {
						return io.NopCloser(strings.NewReader(content)), nil
					}
					return os.Open(path)
				},
				ImportPaths: importPath,
			},
		),
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	files, err := compiler.Compile(context.TODO(), protoFile)
	if err != nil {
		return "", err
	}

	dp, err := files.AsResolver().FindDescriptorByName(protoreflect.FullName(service))
	if err != nil {
		return "", err
	}

	randFuncMap := map[protoreflect.Kind]func(md protoreflect.FieldDescriptor) any{}
	randFuncMap[protoreflect.Int32Kind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		}
		return rand.Intn(100)
	}
	randFuncMap[protoreflect.Uint32Kind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		}
		return rand.Intn(100)
	}
	randFuncMap[protoreflect.Int64Kind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		}
		return rand.Intn(100)
	}
	randFuncMap[protoreflect.Uint64Kind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		}
		return rand.Intn(100)
	}
	randFuncMap[protoreflect.FloatKind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []float32{rand.Float32(), rand.Float32(), rand.Float32()}
		}
		return rand.Float32()
	}
	randFuncMap[protoreflect.DoubleKind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []float64{rand.Float64(), rand.Float64(), rand.Float64()}
		}
		return rand.Float64()
	}
	randFuncMap[protoreflect.BoolKind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []bool{true, false, true}
		}
		return true
	}
	randFuncMap[protoreflect.StringKind] = func(md protoreflect.FieldDescriptor) any {
		if md.IsList() {
			return []string{"xxx", "yyy", "zzz"}
		}
		return "xxx"
	}
	randFuncMap[protoreflect.EnumKind] = func(md protoreflect.FieldDescriptor) any {
		enums := md.Enum().Values()
		return enums.Get(rand.Intn(enums.Len())).Index()
	}
	randFuncMap[protoreflect.MessageKind] = func(md protoreflect.FieldDescriptor) any {
		result := map[string]any{}
		if md.IsMap() {
			key := randFuncMap[md.MapKey().Kind()](md.MapKey())
			if strKey, ok := key.(string); ok {
				result[strKey] = randFuncMap[md.MapValue().Kind()](md.MapValue())
			}
		} else if md.IsList() {
			child := md.Message().Fields().Get(0)
			randFunc := randFuncMap[child.Kind()]
			result[md.JSONName()] = []interface{}{randFunc(child)}
		} else {
			for i := 0; i < md.Message().Fields().Len(); i++ {
				field := md.Message().Fields().Get(i)
				randFunc := randFuncMap[field.Kind()]
				if randFunc != nil {
					result[field.JSONName()] = randFunc(field)
				}
			}
		}
		return result
	}

	data := map[string]any{}
	abc := dp.(protoreflect.MethodDescriptor)
	for i := 0; i < abc.Input().Fields().Len(); i++ {
		field := abc.Input().Fields().Get(i)
		randFunc := randFuncMap[field.Kind()]
		if randFunc != nil {
			data[string(field.Name())] = randFunc(field)
		}
	}

	var result []byte
	result, err = json.Marshal(data)
	if err == nil {
		resultJSON = string(result)
	}
	return
}

func init() {
	RegisterCodeGenerator("gRPCPayload", NewGrpcPayloadGenerator())
}
