package controller

import (
	rtmp "github.com/zhangpeihao/gortmp"
	"log"
	"github.com/instrumentisto/go-rtmp-bot/model"
)

// RTMP clients event handler.
// The implementation of rtmp OutboundHandler
type RTMPHandler struct {
	ID      string
	Handler *AppHandler
}

// Handles changing status of rtmp connection.
// Just calls Application handler onSignal with status signal.
//
// param:  Reference to RTMP connection.
func (h *RTMPHandler) OnStatus(conn rtmp.OutboundConn) {
	var err error
	status, err := conn.Status()
	if err != nil {
		log.Panicf("can not read status: %s", err.Error())
	}
	signal := model.NewSignal(model.STATUS, h.ID)
	signal.Data = status
	h.Handler.OnSignal(signal)
}

// Handles close RTMP connection.
// Just calls Application handler onSignal with closed signal.
//
// params:  Reference to rtmp connection.
func (h *RTMPHandler) OnClosed(conn rtmp.Conn) {
	signal := model.NewSignal(model.CLOSED, h.ID)
	h.Handler.OnSignal(signal)
}

// Handles receiving of any RTMP message.
// Just calls Application handler onSignal with play stream signal.
//
// params:  Reference to RTMP connection;
//          Any RTMP message (in this case - video frame).
func (h *RTMPHandler) OnReceived(conn rtmp.Conn, message *rtmp.Message) {
	signal := model.NewSignal(model.PLAY_STREAM, h.ID)
	signal.Data = message
	h.Handler.OnSignal(signal)
}

// Handles received rtmp command.
// This method implements OutboundHandler interface only.
//
// params:  Reference to RTMP connection;
//          Any RTMP command.
func (h *RTMPHandler) OnReceivedRtmpCommand(
	conn rtmp.Conn, command *rtmp.Command) {
}

// Handles stream creation.
// Just calls Application handler onSignal with stream create signal.
//
// params:  Reference to RTMP connection;
//          Reference to RTMP stream instance.

func (h *RTMPHandler) OnStreamCreated(
	conn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	signal := model.NewSignal(model.STREAM_CREATE, h.ID)
	signal.Data = stream
	h.Handler.OnSignal(signal)
}

// Handles play start.
// This method implements OutboundHandler interface only.
//
// params: Reference to RTMP stream instance.
func (h *RTMPHandler) OnPlayStart(stream rtmp.OutboundStream) {
	// Does nothing.
}

// Handles start of stream publishing.
// Just calls Application handler onSignal with publish start signal.
//
// params: Reference to RTMP stream instance.
func (h *RTMPHandler) OnPublishStart(stream rtmp.OutboundStream) {
	signal := model.NewSignal(model.PUBLISH_START, h.ID)
	signal.Data = stream
	h.Handler.OnSignal(signal)
}
