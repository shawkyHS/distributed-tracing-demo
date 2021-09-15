// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redisreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

func TestRedisRunnable(t *testing.T) {
	consumer := new(consumertest.MetricsSink)
	logger, _ := zap.NewDevelopment()
	runner := newRedisRunnable(context.Background(), config.NewID(typeStr), newFakeClient(), consumer, logger)
	err := runner.Setup()
	require.Nil(t, err)
	err = runner.Run()
	require.Nil(t, err)
	// + 6 because there are two keyspace entries each of which has three metrics
	assert.Equal(t, len(getDefaultRedisMetrics())+6, consumer.DataPointCount())
	md := consumer.AllMetrics()[0]
	rm := md.ResourceMetrics().At(0)
	ilm := rm.InstrumentationLibraryMetrics().At(0)
	il := ilm.InstrumentationLibrary()
	assert.Equal(t, "otelcol/redis", il.Name())

}
