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
