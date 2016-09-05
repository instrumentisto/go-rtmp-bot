package rtmp_bot

import (
	rtmp "github.com/zhangpeihao/gortmp"

	"github.com/instrumentisto/go-rtmp-bot/model"
)

// RTMP client interface (publisher or player).
type IRTMPClient interface {
	SetStatus(status uint)                    // Sets RTMP connection status.
	SetStream(stream rtmp.OutboundStream)     // Sets reference to RTMP stream.
	PublishStream(stream rtmp.OutboundStream) //Publishes stream.
	PlayStream(message *rtmp.Message)         // Plays RTMP stream.
	Stop()                                    // Stops RTMP client.
	GetID() string                            // Returns RTMP client ID.
	Run()                                     // Runs RTMP connection.
	GetStreamKey() string                     // Returns RTMP stream key.
	GetStat() *model.StatItem                 // Returns statistic instance.
	AddFrame(frame *model.FlvFrame)           // Adds rtmp frame
	UpdateStat()                              // Updates statistics.
}
