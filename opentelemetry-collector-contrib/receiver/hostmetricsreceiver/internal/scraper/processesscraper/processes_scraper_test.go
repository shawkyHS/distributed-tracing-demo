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

package processesscraper

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/load"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/model/pdata"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/processesscraper/internal/metadata"
)

func TestScrape(t *testing.T) {
	type testCase struct {
		name              string
		miscFunc          func() (*load.MiscStat, error)
		expectedStartTime pdata.Timestamp
		expectedErr       string
	}

	testCases := []testCase{
		{
			name: "Standard",
		},
		{
			name:        "Error",
			miscFunc:    func() (*load.MiscStat, error) { return nil, errors.New("err1") },
			expectedErr: "err1",
		},
		{
			name:              "Validate Start Time",
			expectedStartTime: 100 * 1e9,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			expectProcessesCountMetric := runtime.GOOS == "linux" || runtime.GOOS == "openbsd" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd"
			expectProcessesCreatedMetric := runtime.GOOS == "linux" || runtime.GOOS == "openbsd"

			scraper := newProcessesScraper(context.Background(), &Config{})
			if test.miscFunc != nil {
				scraper.misc = test.miscFunc
			}

			err := scraper.start(context.Background(), componenttest.NewNopHost())
			require.NoError(t, err, "Failed to initialize processes scraper: %v", err)
			if test.expectedStartTime != 0 {
				scraper.startTime = test.expectedStartTime
			}

			metrics, err := scraper.scrape(context.Background())

			expectedMetricCount := 0
			if expectProcessesCountMetric {
				expectedMetricCount++
			}
			if expectProcessesCreatedMetric {
				expectedMetricCount++
			}

			if (expectProcessesCountMetric || expectProcessesCreatedMetric) && test.expectedErr != "" {
				assert.EqualError(t, err, test.expectedErr)

				isPartial := scrapererror.IsPartialScrapeError(err)
				assert.True(t, isPartial)
				if isPartial {
					assert.Equal(t, expectedMetricCount, err.(scrapererror.PartialScrapeError).Failed)
				}

				return
			}
			require.NoError(t, err, "Failed to scrape metrics: %v", err)

			assert.Equal(t, expectedMetricCount, metrics.Len())

			if expectProcessesCountMetric {
				assertProcessesCountMetricValid(t, metrics.At(0), test.expectedStartTime)
			}
			if expectProcessesCreatedMetric {
				assertProcessesCreatedMetricValid(t, metrics.At(1), test.expectedStartTime)
			}

			internal.AssertSameTimeStampForAllMetrics(t, metrics)
		})
	}
}

func assertProcessesCountMetricValid(t *testing.T, metric pdata.Metric, startTime pdata.Timestamp) {
	internal.AssertDescriptorEqual(t, metadata.Metrics.SystemProcessesCount.New(), metric)
	if startTime != 0 {
		internal.AssertSumMetricStartTimeEquals(t, metric, startTime)
	}
	assert.Equal(t, 2, metric.Sum().DataPoints().Len())
	internal.AssertSumMetricHasAttributeValue(t, metric, 0, "status", pdata.NewAttributeValueString(metadata.LabelStatus.Running))
	internal.AssertSumMetricHasAttributeValue(t, metric, 1, "status", pdata.NewAttributeValueString(metadata.LabelStatus.Blocked))
}

func assertProcessesCreatedMetricValid(t *testing.T, metric pdata.Metric, startTime pdata.Timestamp) {
	if startTime != 0 {
		internal.AssertSumMetricStartTimeEquals(t, metric, startTime)
	}
	internal.AssertDescriptorEqual(t, metadata.Metrics.SystemProcessesCreated.New(), metric)
	assert.Equal(t, 1, metric.Sum().DataPoints().Len())
	assert.Equal(t, 0, metric.Sum().DataPoints().At(0).Attributes().Len())
}
