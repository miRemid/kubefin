/*
Copyright 2023 The KubeFin Authors

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

package ack

type NodeSpecIno struct {
	InstanceTypeId string `json:"instanceTypeId"`
	CPUCoreCount   string `json:"cpuCoreCount"`
	MemorySize     string `json:"memorySize"`
}

type NodeSpecQueryResultInstanceType struct {
	InstanceType []NodeSpecIno `json:"instance_type"`
}

type NodeSpecQueryResultComponent struct {
	InstanceType NodeSpecQueryResultInstanceType `json:"instance_type"`
}

type NodeSpecQueryResultData struct {
	Components NodeSpecQueryResultComponent `json:"components"`
}

type NodeSpecQueryResult struct {
	Data NodeSpecQueryResultData `json:"data"`
}

type TenantCalculatorResultOrder struct {
	TradeAmount float64 `json:"tradeAmount"`
}

type TenantCalculatorResultData struct {
	Order TenantCalculatorResultOrder `json:"order"`
}

type TenantCalculatorResult struct {
	Data TenantCalculatorResultData `json:"data"`
}

type TenantCalculatorInstanceProperty struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}

type TenantCalculatorComponent struct {
	ComponentCode    string                             `json:"componentCode"`
	InstanceProperty []TenantCalculatorInstanceProperty `json:"instanceProperty"`
}

type TenantCalculatorConfiguration struct {
	CommodityCode   string                      `json:"commodityCode"`
	SpecCode        string                      `json:"specCode"`
	ChargeType      string                      `json:"chargeType"`
	OrderType       string                      `json:"orderType"`
	Quantity        int                         `json:"quantity"`
	Duration        int                         `json:"duration"`
	PricingCycle    string                      `json:"pricingCycle"`
	Components      []TenantCalculatorComponent `json:"components"`
	UseTimeUnit     string                      `json:"useTimeUnit"`
	UseTimeQuantity int                         `json:"useTimeQuantity"`
}

type TenantCalculator struct {
	Tenant         string                          `json:"tenant"`
	Configurations []TenantCalculatorConfiguration `json:"configurations"`
}
