// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metricsx

import (
	"testing"
	"time"

	"github.com/apache/beam/sdks/v2/go/pkg/beam/core/metrics"
	pipepb "github.com/apache/beam/sdks/v2/go/pkg/beam/model/pipeline_v1"
	"github.com/google/go-cmp/cmp"
)

func TestFromMonitoringInfos_Counters(t *testing.T) {
	var value int64 = 15
	want := metrics.CounterResult{
		Attempted: 15,
		Committed: 0,
		Key: metrics.StepKey{
			Step:      "main.customDoFn",
			Name:      "customCounter",
			Namespace: "customDoFn",
		}}

	payload, err := Int64Counter(value)
	if err != nil {
		t.Fatalf("Failed to encode Int64Counter: %v", err)
	}

	labels := map[string]string{
		"PTRANSFORM": "main.customDoFn",
		"NAMESPACE":  "customDoFn",
		"NAME":       "customCounter",
	}

	mInfo := &pipepb.MonitoringInfo{
		Urn:     UrnToString(UrnUserSumInt64),
		Type:    UrnToType(UrnUserSumInt64),
		Labels:  labels,
		Payload: payload,
	}

	attempted := []*pipepb.MonitoringInfo{mInfo}
	committed := []*pipepb.MonitoringInfo{}
	p := &pipepb.Pipeline{}

	got := FromMonitoringInfos(p, attempted, committed).AllMetrics().Counters()
	size := len(got)
	if size < 1 {
		t.Fatalf("Invalid array's size: got: %v, want: %v", size, 1)
	}
	if d := cmp.Diff(want, got[0]); d != "" {
		t.Fatalf("Invalid counter: got: %v, want: %v, diff(-want,+got):\n %v",
			got[0], want, d)
	}
}

func TestFromMonitoringInfos_Distributions(t *testing.T) {
	var count, sum, min, max int64 = 100, 5, -12, 30

	want := metrics.DistributionResult{
		Attempted: metrics.DistributionValue{
			Count: 100,
			Sum:   5,
			Min:   -12,
			Max:   30,
		},
		Committed: metrics.DistributionValue{},
		Key: metrics.StepKey{
			Step:      "main.customDoFn",
			Name:      "customDist",
			Namespace: "customDoFn",
		}}

	payload, err := Int64Distribution(count, sum, min, max)
	if err != nil {
		t.Fatalf("Failed to encode Int64Distribution: %v", err)
	}

	labels := map[string]string{
		"PTRANSFORM": "main.customDoFn",
		"NAMESPACE":  "customDoFn",
		"NAME":       "customDist",
	}

	mInfo := &pipepb.MonitoringInfo{
		Urn:     UrnToString(UrnUserDistInt64),
		Type:    UrnToType(UrnUserDistInt64),
		Labels:  labels,
		Payload: payload,
	}

	attempted := []*pipepb.MonitoringInfo{mInfo}
	committed := []*pipepb.MonitoringInfo{}
	p := &pipepb.Pipeline{}

	got := FromMonitoringInfos(p, attempted, committed).AllMetrics().Distributions()
	size := len(got)
	if size < 1 {
		t.Fatalf("Invalid array's size: got: %v, want: %v", size, 1)
	}
	if d := cmp.Diff(want, got[0]); d != "" {
		t.Fatalf("Invalid distribution: got: %v, want: %v, diff(-want,+got):\n %v",
			got[0], want, d)
	}
}

func TestFromMonitoringInfos_Gauges(t *testing.T) {
	var value int64 = 100
	loc, _ := time.LoadLocation("Local")
	tm := time.Date(2020, 11, 9, 17, 52, 28, 462*int(time.Millisecond), loc)

	want := metrics.GaugeResult{
		Attempted: metrics.GaugeValue{
			Value:     100,
			Timestamp: tm,
		},
		Committed: metrics.GaugeValue{},
		Key: metrics.StepKey{
			Step:      "main.customDoFn",
			Name:      "customGauge",
			Namespace: "customDoFn",
		}}

	payload, err := Int64Latest(tm, value)
	if err != nil {
		t.Fatalf("Failed to encode Int64Latest: %v", err)
	}

	labels := map[string]string{
		"PTRANSFORM": "main.customDoFn",
		"NAMESPACE":  "customDoFn",
		"NAME":       "customGauge",
	}

	mInfo := &pipepb.MonitoringInfo{
		Urn:     UrnToString(UrnUserLatestMsInt64),
		Type:    UrnToType(UrnUserLatestMsInt64),
		Labels:  labels,
		Payload: payload,
	}

	attempted := []*pipepb.MonitoringInfo{mInfo}
	committed := []*pipepb.MonitoringInfo{}
	p := &pipepb.Pipeline{}

	got := FromMonitoringInfos(p, attempted, committed).AllMetrics().Gauges()
	size := len(got)
	if size < 1 {
		t.Fatalf("Invalid array's size: got: %v, want: %v", size, 1)
	}
	if d := cmp.Diff(want, got[0]); d != "" {
		t.Fatalf("Invalid gauge: got: %v, want: %v, diff(-want,+got):\n %v",
			got[0], want, d)
	}
}
