package rtmp_bot

import (
	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/player"
	"github.com/instrumentisto/go-rtmp-bot/publisher"
	"github.com/zhangpeihao/gortmp"
	"log"
	"strconv"
	"time"
)

// Stress test launcher.
// Creates publisher and players.
// Collects statistics.
// Makes test report.
const (
	log_filename = "stress_test.log"
)

// RTMP media server stress test launcher.
// Starts requested count of publishers and players.
type Launcher struct {
	Data       *model.StartRequest // Stress test requested parameters.
	TestReport *model.Report
	rtmp_path  string
	clients    map[string]IRTMPClient // Map of RTMP clients.
	handler    *controller.AppHandler // Application signals handler.
	stop_chan  chan bool              // Stop channel.
}

// Constructs new stress test launcher.
//
// params: data *model.StartRequest   Stress test requested parameters.
// return new Launcher instance.
func NewLauncher(data *model.StartRequest, report *model.Report, rtmp_file_path string) *Launcher {
	return &Launcher{
		Data:       data,
		TestReport: report,
		rtmp_path:  rtmp_file_path,
		clients:    make(map[string]IRTMPClient),
		stop_chan:  make(chan bool),
		handler: &controller.AppHandler{
			Signal_chan: make(chan *model.Signal),
		},
	}
}

// Starts stress test.
func (l *Launcher) Start() {
	defer l.cleanMap()
	l.cleanMap()
	l.stop_chan = make(chan bool)
	flv_chan := make(chan *model.FlvFrame)
	flv_stream, err := publisher.NewFlvFile(
		l.rtmp_path,
		l.handler)
	if err != nil {
		return
	}
	go l.startStat()
	defer flv_stream.CloseFile()
	for i := 0; i < l.Data.ModelCount; i++ {
		stream_key := "model" + strconv.Itoa(i+1)
		pub := publisher.NewPublisher(
			l.Data.ServerURL, stream_key, l.handler, flv_chan)
		go pub.Run()
		l.clients[pub.GetID()] = pub
	}
	go flv_stream.PlayFile()
	for {
		select {
		case signal, ok := <-l.handler.Signal_chan:
			if ok {
				switch signal.SignalType {
				case model.STATUS:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					client, ok := l.clients[signal.Target]
					if !ok {
						log.Printf("STATUS client not found: %v", signal.Target)
						continue
					}

					client.SetStatus(signal.Data.(uint))
				case model.CLOSED:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					_, ok := l.clients[signal.Target]
					if !ok {
						log.Printf("CLOSED client not found: %v", signal.Target)
						continue
					}
					break
				case model.STREAM_CREATE:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					client, ok := l.clients[signal.Target]
					if !ok {
						log.Printf("STREAM CREATE client not found: %v", signal.Target)
						continue
					}
					client.SetStream(
						signal.Data.(gortmp.OutboundStream))
				case model.PUBLISH_START:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					client, ok := l.clients[signal.Target]
					if !ok {
						log.Printf("PUBLISH START client not found: %v", signal.Target)
						continue
					}

					client.PublishStream(
						signal.Data.(gortmp.OutboundStream))
					l.startClients(client.GetStreamKey())
				case model.PLAY_STREAM:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					client, ok := l.clients[signal.Target]
					if !ok {
						log.Printf("PLAY STREAM client not found: %v", signal.Target)
						continue
					}
					go client.PlayStream(
						signal.Data.(*gortmp.Message))
				case model.ADD_FRAME:
					if l.clients == nil || len(l.clients) == 0 {
						continue
					}
					for _, c := range l.clients {
						if c.GetStat().Role == model.ROLE_PUBLISHER {
							c.AddFrame(signal.Data.(*model.FlvFrame))
						}
					}
				}
			}
		case <-l.stop_chan:
			return
		}
	}
}

// Starts RTMP players.
//
// param: stream_key string   RTMP stream key.
func (l *Launcher) startClients(stream_key string) {
	for i := 0; i < l.Data.ClientCount; i++ {
		player := player.NewPlayer(l.Data.ServerURL, stream_key, l.handler)
		l.clients[player.GetID()] = player
		go player.Run()
	}
}

// Starts statistics.
func (l *Launcher) startStat() {
	defer l.cleanMap()
	for {
		select {
		case <-l.stop_chan:
			return
		case <-time.After(1 * time.Second):
			l.makeStat()
		}
	}
}

// Writes publishers statistic to log file.
func (l *Launcher) makeStat() {
	client_map := make(map[string]*model.StatItem)
	for _, client := range l.clients {
		client.UpdateStat()
		stat_item := client.GetStat()
		if stat_item.Role == model.ROLE_PUBLISHER {
			l.publisherAddPlayers(stat_item)
		}
		client_map[client.GetID()] = stat_item
	}
	l.TestReport.UpdateReport(client_map)
}

// Writes players statistic to log file.
//
// param: item *model.StatItem   Publisher statistic item.
func (l *Launcher) publisherAddPlayers(item *model.StatItem) {
	for _, client := range l.clients {
		if client.GetStreamKey() == item.StreamID {
			if client.GetStat().Role == model.ROLE_PLAYER {
				item.Receivers[client.GetID()] = client.GetStat()
			}
		}
	}
}

// Returns test report.
func (l *Launcher) updateReport() *model.Report {
	publisher_map := make(map[string]*model.StatItem)
	for _, client := range l.clients {
		stat_item := client.GetStat()
		publisher_map[client.GetID()] = stat_item
	}
	l.TestReport.UpdateReport(publisher_map)
	return l.TestReport
}

// Stops stress test.
func (l *Launcher) Stop() {
	for _, client := range l.clients {
		if client != nil {
			go client.Stop()
		}
	}
	l.stop_chan <- true
}

// Cleans RTMP clients map.
func (l *Launcher) cleanMap() {
	l.clients = make(map[string]IRTMPClient)
}

// Check any panic.
func (l *Launcher) onClose() {
	if r := recover(); r != nil {
		log.Printf("RECOVER on Launcer %s", r)
		l.stop_chan <- true
	}
}
