// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !linux && !windows
// +build !linux,!windows

package memoryscraper

import (
	"github.com/shirou/gopsutil/mem"
	"go.opentelemetry.io/collector/model/pdata"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/memoryscraper/internal/metadata"
)

const memStatesLen = 3

func appendMemoryUsageStateDataPoints(idps pdata.NumberDataPointSlice, now pdata.Timestamp, memInfo *mem.VirtualMemoryStat) {
	initializeMemoryUsageDataPoint(idps.AppendEmpty(), now, metadata.LabelState.Used, int64(memInfo.Used))
	initializeMemoryUsageDataPoint(idps.AppendEmpty(), now, metadata.LabelState.Free, int64(memInfo.Free))
	initializeMemoryUsageDataPoint(idps.AppendEmpty(), now, metadata.LabelState.Inactive, int64(memInfo.Inactive))
}
