package model

// Signal types constants.
const (
	STATUS        string = "status"
	CLOSED        string = "closed"
	STREAM_CREATE string = "stream_create"
	PUBLISH_START string = "publish_start"
	PLAY_STREAM   string = "play_stream"
	ADD_FRAME     string = "add_frame"
)

// Relation of signals model.
type Signal struct {
	SignalType string      // Type of signal
	Target     string      // Signal target identifier.
	Data       interface{} // Any signal data.
}

// Construct new Signal instance.
//
// params: Signal type   string
//         Signal target string
// returns new Signal instance.
func NewSignal(signal_type string, target string) *Signal {
	return &Signal{
		SignalType: signal_type,
		Target:     target,
	}
}
