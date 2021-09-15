// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datareceivers

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/testbed/testbed"
)

// SplunkHECDataReceiver implements Splunk HEC format receiver.
type SplunkHECDataReceiver struct {
	testbed.DataReceiverBase
	receiver component.LogsReceiver
}

// Ensure SplunkHECDataReceiver implements LogDataSender.
var _ testbed.DataReceiver = (*SplunkHECDataReceiver)(nil)

// NewSplunkHECDataReceiver creates a new SplunkHECDataReceiver that will listen on the
// specified port after Start is called.
func NewSplunkHECDataReceiver(port int) *SplunkHECDataReceiver {
	return &SplunkHECDataReceiver{DataReceiverBase: testbed.DataReceiverBase{Port: port}}
}

// Start the receiver.
func (sr *SplunkHECDataReceiver) Start(_ consumer.Traces, _ consumer.Metrics, lc consumer.Logs) error {
	config := splunkhecreceiver.Config{
		HTTPServerSettings: confighttp.HTTPServerSettings{
			Endpoint: fmt.Sprintf("localhost:%d", sr.Port),
		},
	}
	var err error
	f := splunkhecreceiver.NewFactory()
	sr.receiver, err = f.CreateLogsReceiver(context.Background(), component.ReceiverCreateSettings{Logger: zap.L()}, &config, lc)
	if err != nil {
		return err
	}

	return sr.receiver.Start(context.Background(), sr)
}

// Stop the receiver.
func (sr *SplunkHECDataReceiver) Stop() error {
	return sr.receiver.Shutdown(context.Background())
}

// GenConfigYAMLStr returns exporter config for the agent.
func (sr *SplunkHECDataReceiver) GenConfigYAMLStr() string {
	// Note that this generates an exporter config for agent.
	return fmt.Sprintf(`
    splunk_hec:
      endpoint: "http://localhost:%d"
      token: "token"`, sr.Port)
}

// ProtocolName returns protocol name as it is specified in Collector config.
func (sr *SplunkHECDataReceiver) ProtocolName() string {
	return "splunk_hec"
}