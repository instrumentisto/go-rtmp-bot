package controller

import "github.com/instrumentisto/go-rtmp-bot/model"

// Signal communication handler.
type AppHandler struct {
	Signal_chan chan *model.Signal
}

// Writes to signal channel any income signals.
func (h *AppHandler) OnSignal(signal *model.Signal) {
	h.Signal_chan <- signal
}
