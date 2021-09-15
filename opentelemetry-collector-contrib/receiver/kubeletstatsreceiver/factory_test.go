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

//go:build !windows
// +build !windows

// TODO review if tests should succeed on Windows

package kubeletstatsreceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configcheck"
	"go.opentelemetry.io/collector/config/configparser"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	kube "github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	ft := factory.Type()
	require.EqualValues(t, "kubeletstats", ft)
}

func TestValidConfig(t *testing.T) {
	factory := NewFactory()
	err := configcheck.ValidateConfig(factory.CreateDefaultConfig())
	require.NoError(t, err)
}

func TestCreateTracesReceiver(t *testing.T) {
	factory := NewFactory()
	traceReceiver, err := factory.CreateTracesReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		factory.CreateDefaultConfig(),
		nil,
	)
	require.ErrorIs(t, err, componenterror.ErrDataTypeIsNotSupported)
	require.Nil(t, traceReceiver)
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	metricsReceiver, err := factory.CreateMetricsReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		tlsConfig(),
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, metricsReceiver)
}

func TestFactoryInvalidExtraMetadataLabels(t *testing.T) {
	factory := NewFactory()
	cfg := Config{
		ExtraMetadataLabels: []kubelet.MetadataLabel{kubelet.MetadataLabel("invalid-label")},
	}
	metricsReceiver, err := factory.CreateMetricsReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		&cfg,
		consumertest.NewNop(),
	)
	require.Error(t, err)
	require.Equal(t, "label \"invalid-label\" is not supported", err.Error())
	require.Nil(t, metricsReceiver)
}

func TestFactoryBadAuthType(t *testing.T) {
	factory := NewFactory()
	cfg := &Config{
		ClientConfig: kube.ClientConfig{
			APIConfig: k8sconfig.APIConfig{
				AuthType: "foo",
			},
		},
	}
	_, err := factory.CreateMetricsReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		cfg,
		consumertest.NewNop(),
	)
	require.Error(t, err)
}

func TestRestClientErr(t *testing.T) {
	cfg := &Config{
		ClientConfig: kube.ClientConfig{
			APIConfig: k8sconfig.APIConfig{
				AuthType: "tls",
			},
		},
	}
	_, err := restClient(zap.NewNop(), cfg)
	require.Error(t, err)
}

func tlsConfig() *Config {
	return &Config{
		ClientConfig: kube.ClientConfig{
			APIConfig: k8sconfig.APIConfig{
				AuthType: "tls",
			},
			TLSSetting: configtls.TLSSetting{
				CAFile:   "testdata/testcert.crt",
				CertFile: "testdata/testcert.crt",
				KeyFile:  "testdata/testkey.key",
			},
		},
	}
}

func TestCustomUnmarshaller(t *testing.T) {
	type args struct {
		componentParser *configparser.Parser
		intoCfg         *Config
	}
	tests := []struct {
		name                  string
		args                  args
		result                *Config
		mockUnmarshallFailure bool
		configOverride        map[string]interface{}
		wantErr               bool
	}{
		{
			name:    "No config",
			wantErr: false,
		},
		{
			name: "Fail initial unmarshal",
			args: args{
				componentParser: configparser.NewParser(),
			},
			wantErr: true,
		},
		{
			name: "metric_group unset",
			args: args{
				componentParser: configparser.NewParser(),
				intoCfg:         &Config{},
			},
			result: &Config{
				MetricGroupsToCollect: defaultMetricGroups,
			},
		},
		{
			name: "fail to unmarshall metric_groups",
			args: args{
				componentParser: configparser.NewParser(),
				intoCfg:         &Config{},
			},
			mockUnmarshallFailure: true,
			wantErr:               true,
		},
		{
			name: "successfully override metric_group",
			args: args{
				componentParser: configparser.NewParser(),
				intoCfg: &Config{
					CollectionInterval: 10 * time.Second,
				},
			},
			configOverride: map[string]interface{}{
				"metric_groups":       []kubelet.MetricGroup{kubelet.ContainerMetricGroup},
				"collection_interval": 20 * time.Second,
			},
			result: &Config{
				CollectionInterval:    20 * time.Second,
				MetricGroupsToCollect: []kubelet.MetricGroup{kubelet.ContainerMetricGroup},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockUnmarshallFailure {
				// some arbitrary failure.
				tt.args.componentParser.Set(metricGroupsConfig, map[string]string{})
			}

			// Mock some config overrides.
			if tt.configOverride != nil {
				for k, v := range tt.configOverride {
					tt.args.componentParser.Set(k, v)
				}
			}

			if err := tt.args.intoCfg.Unmarshal(tt.args.componentParser); (err != nil) != tt.wantErr {
				t.Errorf("customUnmarshaller() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.result != nil {
				assert.Equal(t, tt.result, tt.args.intoCfg)
			}
		})
	}
}
