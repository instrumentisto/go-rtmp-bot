package prometheus

import (
	"github.com/Zumata/exporttools"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/prometheus/client_golang/prometheus"
)

// Implements Prometheus Exporter interface.
type stressTestExporter struct {
	exporter *exporttools.BaseExporter // Base Prometheus exporter instance.
	report   *model.Report             // Stress test report.
}

// Returns new prometheus stress test exporter instance.
//
// param: report *model.Report   Stress test report instance.
func NewExporter(report *model.Report) *stressTestExporter {
	return &stressTestExporter{
		exporter: exporttools.NewBaseExporter("stress_test"),
		report:   report,
	}
}

// Setups stress test exporter.
func (e *stressTestExporter) Setup() error {
	e.exporter.AddGroup(newMetricsCollector(e.report))
	return nil
}

// Closes stress test exporter.
func (e *stressTestExporter) Close() error {
	return nil
}

// Describes prometheus metrics.
func (e *stressTestExporter) Describe(ch chan<- *prometheus.Desc) {
	exporttools.GenericDescribe(e.exporter, ch)
}

// Collects prometheus metrics.
func (e *stressTestExporter) Collect(ch chan<- prometheus.Metric) {
	exporttools.GenericCollect(e.exporter, ch)
}

// Process metrics.
func (e *stressTestExporter) Process() {
	e.exporter.Process()
}
