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
	"encoding/json"
	"io"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/linuxsuren/api-testing/pkg/runner/monitor"
)

func NewRemoteServer(dockerCli command.Cli) monitor.MonitorServer {
	return &monitorServer{
		dockerCli: dockerCli,
	}
}

type monitorServer struct {
	dockerCli command.Cli
	monitor.UnimplementedMonitorServer
}

func (s *monitorServer) GetResourceUsage(ctx context.Context, target *monitor.Target) (usage *monitor.ResourceUsage, err error) {
	var st types.ContainerStats
	st, err = s.dockerCli.Client().ContainerStatsOneShot(ctx, target.Name)
	if err != nil {
		return
	}

	var data []byte
	if data, err = io.ReadAll(st.Body); err == nil {
		stats := &types.StatsJSON{}
		if err = json.Unmarshal(data, stats); err == nil {
			usage = &monitor.ResourceUsage{
				Memory: stats.MemoryStats.Usage,
			}
		}
	}
	return
}
