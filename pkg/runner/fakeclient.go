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

package runner

import (
	"context"

	"trpc.group/trpc-go/trpc-go/client"
)

type fakeClient struct {
	err error
}

func (f *fakeClient) Invoke(ctx context.Context, reqBody interface{}, rspBody interface{}, opt ...client.Option) (err error) {
	err = f.err
	return
}
