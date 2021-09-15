// Copyright  OpenTelemetry Authors
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

package extractors

import (
	"time"

	cinfo "github.com/google/cadvisor/info/v1"
	"go.uber.org/zap"

	ci "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight"
	awsmetrics "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics"
)

type NetMetricExtractor struct {
	logger         *zap.Logger
	rateCalculator awsmetrics.MetricCalculator
}

func getInterfacesStats(stats *cinfo.ContainerStats) []cinfo.InterfaceStats {
	ifceStats := stats.Network.Interfaces
	if len(ifceStats) == 0 {
		ifceStats = []cinfo.InterfaceStats{stats.Network.InterfaceStats}
	}
	return ifceStats
}

func (n *NetMetricExtractor) HasValue(info *cinfo.ContainerInfo) bool {
	return info.Spec.HasNetwork
}

func (n *NetMetricExtractor) GetValue(info *cinfo.ContainerInfo, _ CPUMemInfoProvider, containerType string) []*CAdvisorMetric {
	var metrics []*CAdvisorMetric

	// Just a protection here, there is no Container level Net metrics
	if (containerType == ci.TypePod && info.Spec.Labels[containerNameLable] != infraContainerName) || containerType == ci.TypeContainer {
		return metrics
	}

	curStats := GetStats(info)
	curIfceStats := getInterfacesStats(curStats)

	// used for aggregation
	var netIfceMetrics []map[string]interface{}

	for _, cur := range curIfceStats {
		mType := getNetMetricType(containerType, n.logger)
		netIfceMetric := make(map[string]interface{})

		infoName := info.Name + containerType + cur.Name //used to identify the network interface
		multiplier := float64(time.Second)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetRxBytes, infoName, float64(cur.RxBytes), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetRxPackets, infoName, float64(cur.RxPackets), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetRxDropped, infoName, float64(cur.RxDropped), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetRxErrors, infoName, float64(cur.RxErrors), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetTxBytes, infoName, float64(cur.TxBytes), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetTxPackets, infoName, float64(cur.TxPackets), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetTxDropped, infoName, float64(cur.TxDropped), curStats.Timestamp, multiplier)
		assignRateValueToField(&n.rateCalculator, netIfceMetric, ci.NetTxErrors, infoName, float64(cur.TxErrors), curStats.Timestamp, multiplier)

		if netIfceMetric[ci.NetRxBytes] != nil && netIfceMetric[ci.NetTxBytes] != nil {
			netIfceMetric[ci.NetTotalBytes] = netIfceMetric[ci.NetRxBytes].(float64) + netIfceMetric[ci.NetTxBytes].(float64)
		}

		netIfceMetrics = append(netIfceMetrics, netIfceMetric)

		metric := newCadvisorMetric(mType, n.logger)
		metric.tags[ci.NetIfce] = cur.Name
		for k, v := range netIfceMetric {
			metric.fields[ci.MetricName(mType, k)] = v
		}

		metrics = append(metrics, metric)
	}

	aggregatedFields := ci.SumFields(netIfceMetrics)
	if len(aggregatedFields) > 0 {
		metric := newCadvisorMetric(containerType, n.logger)
		for k, v := range aggregatedFields {
			metric.fields[ci.MetricName(containerType, k)] = v
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

func NewNetMetricExtractor(logger *zap.Logger) *NetMetricExtractor {
	return &NetMetricExtractor{
		logger:         logger,
		rateCalculator: newFloat64RateCalculator(),
	}
}

func getNetMetricType(containerType string, logger *zap.Logger) string {
	metricType := ""
	switch containerType {
	case ci.TypeNode:
		metricType = ci.TypeNodeNet
	case ci.TypeInstance:
		metricType = ci.TypeInstanceNet
	case ci.TypePod:
		metricType = ci.TypePodNet
	default:
		logger.Warn("net_extractor: net metric extractor is parsing unexpected containerType", zap.String("containerType", containerType))
	}
	return metricType
}
