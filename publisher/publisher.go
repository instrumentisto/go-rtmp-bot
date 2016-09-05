package publisher

import (
	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/utils"
	"github.com/zhangpeihao/goflv"
	rtmp "github.com/zhangpeihao/gortmp"
	"log"
	"time"
)

// RTMP publisher.
// Publishes RTMP stream to media server.
type Publisher struct {
	status             uint                     // RTMP connection status.
	createStreamChan   chan rtmp.OutboundStream // The channel for created RTMP stream instance.
	serverURL          string                   // Media server URL.
	streamID           string                   // Stream key.
	stop_chanel        chan bool                // The channel for stop publisher instance.
	test_handler       *controller.AppHandler   // Application signal handler reference.
	FlvChan            chan *model.FlvFrame     // Test .flv file url for streaming.
	obConn             rtmp.OutboundConn        // RTMP connection reference.
	id                 string                   // RTMP client identifier.
	stat               *model.StatItem          // Statistic item instance.
	start_command_time int64                    // Start command UNIX time.
	startedAt          int64                    // Publish started UNIX time.
	old_frame_count    int64                    // Count of sends video frames.
	published_stream   rtmp.OutboundStream
}

// Constructs new RTMP Publisher instance.
//
// params: Media server URL             string;
//         Stream key                   string;
//         Application signal handler   *controller.AppHandler
//
// return new *Publisher instance.
func NewPublisher(
	url string, stream_key string,
	test_handler *controller.AppHandler,
	flv_chan chan *model.FlvFrame) *Publisher {
	client_id := utils.GetUUID()
	return &Publisher{
		status:             uint(0),
		serverURL:          url,
		streamID:           stream_key,
		stop_chanel:        make(chan bool),
		test_handler:       test_handler,
		id:                 client_id,
		stat:               model.NewStatItem(model.ROLE_PUBLISHER, stream_key, client_id),
		start_command_time: 0,
		startedAt:          0,
		old_frame_count:    0,
		FlvChan:            flv_chan,
	}
}

// Runs publish stream.
func (p *Publisher) Run() {
	p.start_command_time = time.Now().Unix()
	p.createStreamChan = make(chan rtmp.OutboundStream)
	var err error
	testHandler := &controller.RTMPHandler{
		Handler: p.test_handler,
		ID:      p.id,
	}
	p.obConn, err = rtmp.Dial(p.serverURL, testHandler, 100)
	if err != nil {
		log.Printf("publisher dial error %s", err.Error())
		p.stat.Status = model.STATUS_DESCRIPTIONS[6]
		return
	}
	defer p.obConn.Close()
	err = p.obConn.Connect()
	if err != nil {
		log.Printf("publisher connection error %s", err.Error())
		p.stat.Status = model.STATUS_DESCRIPTIONS[6]
		return
	}
	for {
		select {
		case stream := <-p.createStreamChan:
			stream.Attach(testHandler)
			err = stream.Publish(p.streamID, "live")
			if err != nil {
				log.Printf("publisher publish error %s", err.Error())
				p.stat.Status = model.STATUS_DESCRIPTIONS[6]
				return
			}
		case <-p.stop_chanel:
			return
		}
	}
}

//Publishes test .flv file data to existed RTMP stream.
//
// param: stream   rtmp.OutboundStream
func (p *Publisher) PublishStream(stream rtmp.OutboundStream) {
	p.published_stream = stream
	p.startedAt = time.Now().Unix()
}

func (p *Publisher) AddFrame(frame *model.FlvFrame) {
	if p.published_stream == nil {
		return
	}

	if p.status != rtmp.OUTBOUND_CONN_STATUS_CREATE_STREAM_OK {
		return
	}
	switch frame.Header.TagType {
	case flv.AUDIO_TAG:
		if p.stat.AudioBytes == 0 {
			p.stat.AudioStartUpTime = time.Now().Unix() - p.start_command_time
		}
		p.stat.AudioBytes += int64(len(frame.Frame))
	case flv.VIDEO_TAG:
		if p.stat.VideoBytes == 0 {
			p.stat.VideoStartUpTime = time.Now().Unix() - p.start_command_time
		}
		p.stat.VideoBytes += int64(len(frame.Frame))
		p.stat.TotalFrames++
	}

	if err := p.published_stream.PublishData(
		frame.Header.TagType, frame.Frame,
		frame.DeltaTimestamp); err != nil {
		log.Printf("publish data ERROR: %s", err.Error())
		p.SetStatus(rtmp.OUTBOUND_CONN_STATUS_CLOSE)
		p.published_stream.Close()
		p.published_stream = nil
		return
	}
}

// This method implements IRTMPClient interface only.
//
// param: rtmp stream   rtmp.OutboundStream
func (p *Publisher) PlayStream(message *rtmp.Message) {
	// Does nothing
}

// Stops publisher.
func (p *Publisher) Stop() {
	if p.published_stream != nil {
		p.published_stream.Close()
	}
	p.stop_chanel <- true
}

// Sets RTMP connection status.
//
// param: status uint
func (p *Publisher) SetStatus(status uint) {
	p.status = status
	p.stat.Status = model.STATUS_DESCRIPTIONS[p.status]
}

// Sets reference to existed RTMP stream.
//
// param: stream rtmp.OutboundStream
func (p *Publisher) SetStream(stream rtmp.OutboundStream) {
	p.createStreamChan <- stream
}

// Returns RTMP client identifier.
//
// return string.
func (p *Publisher) GetID() string {
	return p.id
}

// Returns RTMP stream key.
//
// return string
func (p *Publisher) GetStreamKey() string {
	return p.streamID
}

// Returns statistic item instance.
//
// return StatItem.
func (p *Publisher) GetStat() *model.StatItem {
	return p.stat
}

// Updates client statistic.
func (p *Publisher) UpdateStat() {
	if p.status == rtmp.OUTBOUND_CONN_STATUS_CREATE_STREAM_OK {
		p.stat.TotalTime = time.Now().Unix() - p.startedAt
		p.stat.FPS = p.stat.TotalFrames - p.old_frame_count
		p.old_frame_count = p.stat.TotalFrames
	}
}

// Check any panic.
func (p *Publisher) onRecover() {
	if r := recover(); r != nil {
		log.Printf("RECOVER on Publisher %s", r)
		p.Stop()
	}
}
