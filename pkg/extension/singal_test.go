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
package extension

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterStopSignal(t *testing.T) {
	var stoppedA bool
	fs := &fakeServer{}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	cancel()
	RegisterStopSignal(ctx, func() {
		stoppedA = true
	}, fs)
	time.Sleep(time.Second * 2)
	assert.True(t, stoppedA)
	assert.True(t, fs.signal)
}

type fakeServer struct {
	signal bool
}

func (s *fakeServer) Stop() {
	s.signal = true
}
