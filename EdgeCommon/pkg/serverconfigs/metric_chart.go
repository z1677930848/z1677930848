// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type MetricChartType = string

const (
	MetricChartTypePie      MetricChartType = "pie"
	MetricChartTypeBar      MetricChartType = "bar"
	MetricChartTypeTimeBar  MetricChartType = "timeBar"
	MetricChartTypeTimeLine MetricChartType = "timeLine"
	MetricChartTypeTable    MetricChartType = "table"
)

func FindAllMetricChartTypes() []*shared.Definition {
	return []*shared.Definition{
		{Name: "Bar", Code: MetricChartTypeBar, Description: "bar chart", Icon: "chart bar"},
		{Name: "Pie", Code: MetricChartTypePie, Description: "pie chart", Icon: "chart pie"},
		{Name: "Time Bar", Code: MetricChartTypeTimeBar, Description: "time based bar chart", Icon: "chart bar"},
		{Name: "Time Line", Code: MetricChartTypeTimeLine, Description: "time based line chart", Icon: "chart line area"},
		{Name: "Table", Code: MetricChartTypeTable, Description: "table view", Icon: "table"},
	}
}

func FindMetricChartTypeName(chartType MetricChartType) string {
	for _, def := range FindAllMetricChartTypes() {
		if def.Code == chartType {
			return def.Name
		}
	}
	return ""
}
