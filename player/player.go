package player

import (
	rtmp "github.com/zhangpeihao/gortmp"
	"log"

	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/model"
	"github.com/instrumentisto/go-rtmp-bot/utils"
	"time"
)

// RTMP player.
// Plays RTMP stream from media server.
type Player struct {
	status             uint                     // RTMP connection status.
	createStreamChan   chan rtmp.OutboundStream // The channel for created RTMP stream instance.
	id                 string                   // RTMP client identifier.
	serverURL          string                   // Media server URL.
	streamID           string                   // Stream key.
	stop_chanel        chan bool                // The channel for stop player instance.
	test_handler       *controller.AppHandler   // Application signal handler reference.
	obConn             rtmp.OutboundConn        // RTMP connection reference.
	stat               *model.StatItem          // Statistic item instance.
	start_command_time int64                    // Start command UNIX time.
	startedAt          int64                    // Player started UNIX time.
	old_frame_count    int64                    // Count of receiving video frames.
}

// Constructs new RTMP player instance.
//
// params: Media-server URL                       string
//         Stream key                             string
//         Application signal handler reference   *controller.AppHandler
//
// returns: new instance of Player
func NewPlayer(
	url string, stream_key string, test_handler *controller.AppHandler) *Player {
	client_id := utils.GetUUID()
	return &Player{
		status:          uint(0),
		serverURL:       url,
		streamID:        stream_key,
		stop_chanel:     make(chan bool),
		test_handler:    test_handler,
		id:              client_id,
		startedAt:       0,
		stat:            model.NewStatItem(model.ROLE_PLAYER, stream_key, client_id),
		old_frame_count: 0,
	}
}

// Runs RTMP player.
func (p *Player) Run() {
	defer p.onRecover()
	p.createStreamChan = make(chan rtmp.OutboundStream)
	p.start_command_time = time.Now().Unix()
	testHandler := &controller.RTMPHandler{
		Handler: p.test_handler,
		ID:      p.id,
	}
	var err error
	p.obConn, err = rtmp.Dial(p.serverURL, testHandler, 100)

	if err != nil {
		p.stat.Status = model.STATUS_DESCRIPTIONS[6]
		log.Printf("Player DIAL error: %s", err.Error())
		return
	}
	defer p.obConn.Close()
	err = p.obConn.Connect()

	if err != nil {
		log.Printf("Player CONNECTION error: %s", err.Error())
		p.stat.Status = model.STATUS_DESCRIPTIONS[6]
		return
	}
	for {
		select {
		case stream := <-p.createStreamChan:
			// Play
			err = stream.Play(p.streamID, nil, nil, nil)
			if err != nil {
				p.stat.Status = model.STATUS_DESCRIPTIONS[6]
				log.Printf("Player PLAY error: %s", err.Error())
				return
			}
		case stop_command := <-p.stop_chanel:
			if stop_command {
				return
			}
		}
	}
}

// Process of RTMP message.
//
// param: message   rtmp.Message.
func (p *Player) PlayStream(message *rtmp.Message) {
	switch message.Type {
	case rtmp.VIDEO_TYPE:

		if p.stat.VideoBytes == 0 {
			p.stat.VideoStartUpTime = time.Now().Unix() - p.start_command_time
		}

		if p.startedAt == 0 && p.stat.VideoBytes > 0 {
			p.startedAt = time.Now().Unix()
		}
		p.stat.VideoBytes += int64(message.Buf.Len())
		p.stat.TotalFrames++
	case rtmp.AUDIO_TYPE:

		if p.stat.AudioBytes == 0 {
			p.stat.AudioStartUpTime = time.Now().Unix() - p.start_command_time
		}

		p.stat.AudioBytes += int64(message.Buf.Len())
	}
}

// This method implements IRTMPClient interface only.
//
// param: rtmp stream   rtmp.OutboundStream
func (p *Player) PublishStream(stream rtmp.OutboundStream) {
	// Does nothing!
}
func (p *Player) AddFrame(frame *model.FlvFrame) {
	// Does nothing!
}

// Stops player.
func (p *Player) Stop() {
	if p.status != rtmp.INBOUND_CONN_STATUS_CLOSE {
		p.stop_chanel <- true
	}

}

// Sets RTMP connection status.
//
// param: status uint.
func (p *Player) SetStatus(status uint) {
	p.status = status
	p.stat.Status = model.STATUS_DESCRIPTIONS[p.status]
}

// Sets created RTMP stream reference.
func (p *Player) SetStream(stream rtmp.OutboundStream) {
	p.createStreamChan <- stream
}

// Returns RTMP stream key.
//
// return string.
func (p *Player) GetStreamKey() string {
	return p.streamID
}

// Returns the RTMP client identifier.
//
// return string.
func (p *Player) GetID() string {
	return p.id
}

// Returns statistic item instance.
//
// return StatItem.
func (p *Player) GetStat() *model.StatItem {
	return p.stat
}

// Updates client statistic.
func (p *Player) UpdateStat() {

	if p.stat.VideoBytes > 0 {
		p.stat.TotalTime = time.Now().Unix() - p.startedAt
	}

	if p.stat.TotalFrames != p.old_frame_count {
		p.stat.FPS = p.stat.TotalFrames - p.old_frame_count
		p.old_frame_count = p.stat.TotalFrames
	}
}

// Check any panic.
func (p *Player) onRecover() {
	if r := recover(); r != nil {
		log.Printf("RECOVER on Player %s", r)
		p.Stop()
	}
}
