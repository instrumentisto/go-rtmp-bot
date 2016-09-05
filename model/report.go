package model

import (
	"log"
	"time"
)

// Stress test report.
type Report struct {
	TestId                 string // Test ID.
	StartTime              int64  // Test starts time.
	TotalTime              int64  // Total test time.
	TotalClients           int    // Total RTMP clients count.
	RequestedModelsCount   int    // Requested publishers count.
	RequestedClientsCount  int    // Requested players count.
	ConnectedModelsCount   int64  // Connected publishers count.
	ConnectedClientsCount  int64  // Connected players count.
	ConnectedModelCountLag int64  // Lag with requested and connected
	// publishers.
	ConnectedClientCountLag int64 // Lag with requested and connected
	// players.
	AverageModelFPS           int64 // Average publisher FPS value.
	AverageClientFPS          int64 // Average player FPS value.
	AverageAudioBytesSends    int64 // Average audio bytes published.
	AverageVideoBytesSends    int64 // Average video bytes published.
	AverageAudioBytesReceived int64 // Average audio bytes received.
	AverageVideoBytesReceived int64 // Average video bytes received.
	TotalVideoPublished       int64 // Video data published in bytes.
	TotalVideoPlayed          int64 // Video data received in bytes.
}

// Returns new stress test report instance.
func NewReport() *Report {
	report := &Report{}
	report.ResetReport("", 0, 0)
	return report
}

// Resets report data.
//
// params: test_id      string    Identifier of current test.
//         model_count  int       Count of models.
//         client_count int       Count of clients.
func (r *Report) ResetReport(test_id string, model_count int, client_count int) {
	log.Printf("RESET report: model_count %d client count %d", model_count, client_count)
	r.TestId = test_id
	r.StartTime = time.Now().Unix()
	r.TotalTime = time.Unix(0, 0).Unix()
	r.TotalClients = model_count*client_count + model_count
	r.RequestedModelsCount = model_count
	r.RequestedClientsCount = client_count * model_count
	r.ConnectedModelsCount = 0
	r.ConnectedClientsCount = 0
	r.ConnectedModelCountLag = 0
	r.ConnectedClientCountLag = 0
	r.AverageModelFPS = 0
	r.AverageClientFPS = 0
	r.AverageAudioBytesSends = 0
	r.AverageVideoBytesSends = 0
	r.AverageAudioBytesReceived = 0
	r.AverageVideoBytesReceived = 0
	r.TotalVideoPublished = r.TotalTime
	r.TotalVideoPlayed = r.TotalTime
}

// Updates stress test report.
//
// param: clients map   RTMP clients statistic map.
func (r *Report) UpdateReport(clients map[string]*StatItem) {
	var total_model_fps int64 = 0
	var total_client_fps int64 = 0
	var publisher_video_start_delay_sum int64 = 0
	var player_video_start_delay_sum int64 = 0
	var publisher_audio_start_delay_sum int64 = 0
	var player_audio_start_delay_sum int64 = 0
	var video_bytes_sends int64 = 0
	var audio_bytes_sends int64 = 0
	var video_bytes_received int64 = 0
	var audio_bytes_received int64 = 0
	var published_total_time int64 = 0
	var played_total_time int64 = 0
	r.ConnectedModelsCount = 0
	r.ConnectedClientsCount = 0
	for _, client := range clients {
		if client.Role == ROLE_PUBLISHER &&
			client.Status == STATUS_DESCRIPTIONS[5] && client.FPS > 0 {
			r.ConnectedModelsCount += 1
			total_model_fps += client.FPS
			publisher_video_start_delay_sum += client.VideoStartUpTime
			publisher_audio_start_delay_sum += client.AudioStartUpTime
			video_bytes_sends += client.VideoBytes
			audio_bytes_sends += client.AudioBytes
			published_total_time += client.TotalTime
		}
		if client.Role == ROLE_PLAYER &&
			client.Status == STATUS_DESCRIPTIONS[5] && client.FPS > 0 {
			r.ConnectedClientsCount += 1
			total_client_fps += client.FPS
			player_video_start_delay_sum += client.VideoStartUpTime
			player_audio_start_delay_sum += client.AudioStartUpTime
			video_bytes_received += client.VideoBytes
			audio_bytes_received += client.AudioBytes
			played_total_time += client.TotalTime
		}
		r.ConnectedModelCountLag = int64(
			r.RequestedModelsCount) - r.ConnectedModelsCount
		r.ConnectedClientCountLag = int64(
			r.RequestedClientsCount) - r.ConnectedClientsCount
		if r.ConnectedModelsCount != 0 {
			r.AverageModelFPS = total_model_fps / int64(
				r.ConnectedModelsCount)
			r.AverageAudioBytesSends = audio_bytes_sends / r.ConnectedModelsCount / 1024
			r.AverageVideoBytesSends = video_bytes_sends / r.ConnectedModelsCount / 1024
			r.TotalVideoPublished = published_total_time / r.ConnectedModelsCount
		}
		if r.ConnectedClientsCount != 0 {
			r.AverageClientFPS = total_client_fps / r.ConnectedClientsCount
			r.AverageAudioBytesReceived = audio_bytes_received / r.ConnectedClientsCount / 1024
			r.AverageVideoBytesReceived = video_bytes_received / r.ConnectedClientsCount / 1024
			r.TotalVideoPlayed = played_total_time / r.ConnectedClientsCount
		}
		r.TotalTime = time.Now().Unix() - r.StartTime
	}
}
