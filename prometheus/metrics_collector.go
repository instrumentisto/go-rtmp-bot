package prometheus

import (
	"github.com/Zumata/exporttools"
	"log"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/utils"
)

// Collector of test metrics.
type metricsCollector struct {
	report *model.Report // Stress test report instance.
}

// Returns new instance of Metrics collector
//
// param: report *model.Report   Instance of the test report.
func newMetricsCollector(report *model.Report) *metricsCollector {
	return &metricsCollector{
		report: report,
	}
}

// Returns array of prometheus metrics or error instance.
func (c *metricsCollector) Collect() ([]*exporttools.Metric, error) {
	log.Println("Collector COLLECT METRICS: %d", c.report.TotalClients)
	metrics := []*exporttools.Metric{
		&exporttools.Metric{
			Name:        "model_connected",
			Type:        exporttools.Gauge,
			Value:       c.report.ConnectedModelsCount,
			Description: "Count of connected models",
		},
		&exporttools.Metric{
			Name:        "clients_connected",
			Type:        exporttools.Gauge,
			Value:       c.report.ConnectedClientsCount,
			Description: "Count of connected clients",
		},
		&exporttools.Metric{
			Name:        "failure_models",
			Type:        exporttools.Gauge,
			Value:       c.report.ConnectedModelCountLag,
			Description: "Count of failures publisher connections",
		},
		&exporttools.Metric{
			Name:        "client_failures",
			Type:        exporttools.Gauge,
			Value:       c.report.ConnectedClientCountLag,
			Description: "Count of failures clients connections",
		},
		&exporttools.Metric{
			Name:        "total_time",
			Type:        exporttools.Gauge,
			Value:       c.report.TotalTime,
			Description: "Stress test total time",
		},
		&exporttools.Metric{
			Name:        "total_clients",
			Type:        exporttools.Gauge,
			Value:       utils.Num64(c.report.TotalClients),
			Description: "Total clients count",
		},
		&exporttools.Metric{
			Name:        "total_model_fps",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageModelFPS,
			Description: "Average model fps",
		},
		&exporttools.Metric{
			Name:        "total_client_fps",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageClientFPS,
			Description: "Average client fps",
		},
		&exporttools.Metric{
			Name:        "audio_bytes_sends",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageAudioBytesSends,
			Description: "Average audio bytes sends",
		},
		&exporttools.Metric{
			Name:        "video_bytes_sends",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageVideoBytesSends,
			Description: "Average video bytes sends",
		},
		&exporttools.Metric{
			Name:        "audio_bytes_received",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageAudioBytesReceived,
			Description: "Average audio bytes received",
		},
		&exporttools.Metric{
			Name:        "video_bytes_received",
			Type:        exporttools.Gauge,
			Value:       c.report.AverageVideoBytesReceived,
			Description: "Average video bytes received",
		},
		&exporttools.Metric{
			Name:        "average_video_time_published",
			Type:        exporttools.Gauge,
			Value:       c.report.TotalVideoPublished,
			Description: "Average video time published",
		},
		&exporttools.Metric{
			Name:        "average_video_time_received",
			Type:        exporttools.Gauge,
			Value:       c.report.TotalVideoPlayed,
			Description: "Average video time received",
		},
	}
	return metrics, nil
}
