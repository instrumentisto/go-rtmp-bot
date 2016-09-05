package model

// Map of RTMP connection status descriptions.
var STATUS_DESCRIPTIONS = map[uint]string{
	0: "Connection closed",
	1: "Connection handshake OK",
	2: "Connect",
	3: "Connect OK",
	4: "Create stream",
	5: "Create stream OK",
	6: "Connection error",
}

// RTMP client roles.
const (
	ROLE_PUBLISHER string = "role_publisher" // Role publisher.
	ROLE_PLAYER    string = "role_player"    // Role player.
)

// Item of Stream statistics.
type StatItem struct {
	Role             string               // Role of RTMP client.
	Status           string               // RTMP connection status.
	StreamID         string               // Stream key.
	ClientID         string               // RTMP client  identifier.
	AudioBytes       int64                // Processed audio bytes.
	VideoBytes       int64                // Processed video bytes.
	VideoStartUpTime int64                // Video publish/play startup time in second.
	AudioStartUpTime int64                // Audio publish/play startup time in second.
	TotalTime        int64                // Total publish/play time in second.
	FPS              int64                // Frames per second.
	Receivers        map[string]*StatItem // Stream receivers map (for publisher only).
	TotalFrames      int64                // Total count of processed RTMP frames.
}

// Constructs new StatItem instance.
//
// params:   Role of statistic target   string;
//           Stream key                 string;
//           RTMP client identifier     string.
//
// returns:  new StatItem instance.
//
func NewStatItem(role string, stream_id string, client_id string) *StatItem {
	return &StatItem{
		Role:             role,
		Status:           "",
		StreamID:         stream_id,
		ClientID:         client_id,
		AudioBytes:       0,
		VideoBytes:       0,
		VideoStartUpTime: 0,
		AudioStartUpTime: 0,
		TotalTime:        0,
		FPS:              0,
		Receivers:        make(map[string]*StatItem),
		TotalFrames:      0,
	}
}
