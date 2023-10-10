/*
Copyright 2022 The KubeFin Authors

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

package query

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog/v2"

	"github.com/kubefin/kubefin/pkg/values"
	"github.com/prometheus/common/model"
)

const (
	instantQueryBaseUrl = "/api/v1/query"
	rangeQueryBaseUrl   = "/api/v1/query_range"
)

type PromqlInstantMessageType struct {
	Data struct {
		ResultType string          `json:"resultType"`
		Result     []*model.Sample `json:"result"`
	} `json:"data"`

	PromqlStatusType
}

type PromqlScalarMessageType struct {
	Data struct {
		ResultType string        `json:"resultType"`
		Result     *model.Scalar `json:"result"`
	} `json:"data"`

	PromqlStatusType
}

type PromqlRangeMessageType struct {
	Data struct {
		ResultType string                `json:"resultType"`
		Result     []*model.SampleStream `json:"result"`
	} `json:"data"`
	PromqlStatusType
}

type PromqlStatusType struct {
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	// Extra field supported by Thanos Querier.
	Warnings []string `json:"warnings"`
}

type PromQueryClient struct {
	endpoint   string
	httpClient *http.Client
	// This is used in multi-tenant storage system
	tenantId string
}

var promQueryClient *PromQueryClient

func InitPromQueryClient(endpoint string) {
	promQueryClient = &PromQueryClient{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: time.Second * 30},
		tenantId:   "",
	}
}

func GetPromQueryClient() *PromQueryClient {
	return promQueryClient
}

func (p *PromQueryClient) WithTenantId(id string) *PromQueryClient {
	clientCopy := *p
	clientCopy.tenantId = id
	return &clientCopy
}

func (p *PromQueryClient) QueryInstant(promql string) ([]*model.Sample, error) {
	respBody, err := p.queryInstant(promql, "")
	if err != nil {
		return nil, err
	}

	message := &PromqlInstantMessageType{}
	if err := json.Unmarshal(respBody, message); err != nil {
		klog.Errorf("Unmarshal error:%v", err)
	}

	return message.Data.Result, nil
}

func (p *PromQueryClient) QueryInstantWithTime(promql string, time int64) ([]*model.Sample, error) {
	timeStr := fmt.Sprintf("%d", time)
	respBody, err := p.queryInstant(promql, timeStr)

	if err != nil {
		return nil, err
	}

	message := &PromqlInstantMessageType{}
	if err := json.Unmarshal(respBody, message); err != nil {
		klog.Errorf("Unmarshal error:%v", err)
		return nil, err
	}

	return message.Data.Result, nil
}

func (p *PromQueryClient) QueryInstantRange(promql string) ([]*model.SampleStream, error) {
	respBody, err := p.queryInstant(promql, "")
	if err != nil {
		return nil, err
	}

	message := &PromqlRangeMessageType{}
	if err := json.Unmarshal(respBody, message); err != nil {
		klog.Errorf("Unmarshal error:%v", err)
	}

	return message.Data.Result, nil
}

func (p *PromQueryClient) queryInstant(promql string, time string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, p.endpoint+instantQueryBaseUrl, nil)
	if err != nil {
		klog.Errorf("Create http request error:%v", err)
		return nil, err
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("query", promql)
	// If it's not set, query from current time
	if time != "" {
		queryParameters.Add("time", time)
	}
	if p.tenantId != "" {
		klog.Infof("Query data with tenant id:%s", p.tenantId)
		req.Header.Add(values.MultiTenantHeader, p.tenantId)
	}
	req.URL.RawQuery = queryParameters.Encode()

	resp, err := p.httpClient.Do(req)
	if err != nil {
		klog.Errorf("Promql query instant error:%v", err)
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Read resp body error:%v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("query failed:%s", string(bodyBytes))
		klog.Errorf("%v", err)
		return nil, err
	}

	return bodyBytes, nil
}

func (p *PromQueryClient) QueryRange(promql string, start, end int64) ([]*model.SampleStream, error) {
	req, err := http.NewRequest(http.MethodGet, p.endpoint+rangeQueryBaseUrl, nil)
	if err != nil {
		klog.Errorf("Create http request error:%v", err)
		return nil, err
	}

	// The returned max point's number is 11000, so we chould choose a right step step seconds
	stepSeconds := (end - start) / 10000
	if stepSeconds < 15 {
		stepSeconds = 15
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("query", promql)
	queryParameters.Add("start", fmt.Sprintf("%d", start))
	queryParameters.Add("end", fmt.Sprintf("%d", end))
	// KubeFin collect metrics in 15s period
	queryParameters.Add("step", fmt.Sprintf("%ds", stepSeconds))
	if p.tenantId != "" {
		klog.V(4).Infof("Query data with tenant id:%s", p.tenantId)
		req.Header.Add(values.MultiTenantHeader, p.tenantId)
	}
	req.URL.RawQuery = queryParameters.Encode()
	resp, err := p.httpClient.Do(req)
	if err != nil {
		klog.Errorf("Promql query range error:%v", err)
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Read resp body error:%v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("query backend error:%s", string(bodyBytes))
		klog.Errorf("%v", err)
		return nil, err
	}

	message := &PromqlRangeMessageType{}
	if err := json.Unmarshal(bodyBytes, message); err != nil {
		klog.Errorf("Unmarshal error:%v", err)
	}

	return message.Data.Result, nil
}

func (p *PromQueryClient) QueryRangeWithStep(promql string, start, end, stepSeconds int64) ([]*model.SampleStream, error) {
	req, err := http.NewRequest(http.MethodGet, p.endpoint+rangeQueryBaseUrl, nil)
	if err != nil {
		klog.Errorf("Create http request error:%v", err)
		return nil, err
	}

	queryParameters := req.URL.Query()
	queryParameters.Add("query", promql)
	queryParameters.Add("start", fmt.Sprintf("%d", start))
	queryParameters.Add("end", fmt.Sprintf("%d", end))
	queryParameters.Add("step", fmt.Sprintf("%ds", stepSeconds))
	if p.tenantId != "" {
		klog.Infof("Query data with tenant id:%s", p.tenantId)
		req.Header.Add(values.MultiTenantHeader, p.tenantId)
	}
	req.URL.RawQuery = queryParameters.Encode()
	resp, err := p.httpClient.Do(req)
	if err != nil {
		klog.Errorf("Promql query range error:%v", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Read resp body error:%v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("query backend error:%s", string(bodyBytes))
		klog.Errorf("%v", err)
		return nil, err
	}

	message := &PromqlRangeMessageType{}
	if err := json.Unmarshal(bodyBytes, message); err != nil {
		klog.Errorf("Unmarshal error:%v", err)
	}

	return message.Data.Result, nil
}
