package prometheus

import (
	"github.com/Zumata/exporttools"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

// Prometheus metrics exporter client.
type ReportExportClient struct {
	listen_address string // Listen address for process metrics response from
	// Prometheus.
	telemetry_path string        // Path to metrics.
	report         *model.Report // Stress test report instance.
}

// Returns new Prometheus client instance.
//
// params: listen_address string          Listen address for process metrics
//                                        response from Prometheus.
//         metrics_path   string          Path to metrics.
//         report         *model.Report   Stress test report instance.
func NewReportExportClient(
	listen_address string,
	metrics_path string,
	report *model.Report) *ReportExportClient {
	return &ReportExportClient{
		listen_address: listen_address,
		telemetry_path: metrics_path,
		report:         report,
	}
}

// Runs prometheus exporter client.
func (c *ReportExportClient) Run() {
	exporter := NewExporter(c.report)
	err := exporttools.Export(exporter)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle(c.telemetry_path, prometheus.Handler())
	http.HandleFunc("/", exporttools.DefaultMetricsHandler(
		"Media server stress test exporter", c.telemetry_path))
	err = http.ListenAndServe(c.listen_address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
