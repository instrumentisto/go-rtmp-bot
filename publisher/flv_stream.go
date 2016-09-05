package publisher

import (
	"github.com/zhangpeihao/goflv"
	"log"
	"time"
	"github.com/instrumentisto/go-rtmp-bot/controller"
	"github.com/instrumentisto/go-rtmp-bot/model"
)

// Opens and read test flv file.
type FlvStream struct {
	FlvFile  *flv.File              // flv file content.
	fileName string                 // flv file name.
	handler  *controller.AppHandler // Application signal handler.
}

// Creates new instance of FlvStream.
//
// params: file_name   string                   Test FLV file name.
//         app_handler *controller.AppHandler   Application signals handler.
func NewFlvFile(
	file_name string,
	app_handler *controller.AppHandler) (*FlvStream, error) {
	file, err := flv.OpenFile(file_name)
	if err != nil {
		return nil, err
	}
	return &FlvStream{
		FlvFile:  file,
		fileName: file_name,
		handler:  app_handler,
	}, nil
}

// Plays the test flv file.
func (s *FlvStream) PlayFile() {
	startTs := uint32(0)
	startAt := time.Now().UnixNano()
	preTs := uint32(0)
	for {
		if s.FlvFile.IsFinished() {
			s.FlvFile.LoopBack()
			startAt = time.Now().UnixNano()
			startTs = uint32(0)
			preTs = uint32(0)
		}
		header, data, err := s.FlvFile.ReadTag()

		if err != nil {
			log.Printf("Read rtmp tag ERROR: %s", err.Error())
			return
		}

		if startTs == uint32(0) {
			startTs = header.Timestamp
		}

		delta_timestamp := uint32(0)

		if header.Timestamp > startTs {
			delta_timestamp = header.Timestamp - startTs
		}
		if delta_timestamp > preTs {
			preTs = delta_timestamp
		}

		frame := &model.FlvFrame{
			Header: header,
			Frame:  data,
		}
		signal := model.NewSignal(model.ADD_FRAME, "flv_stream")
		signal.Data = frame
		s.handler.OnSignal(signal)
		delta2 := uint32((time.Now().UnixNano() - startAt) / 1000000)

		if delta_timestamp > delta2+100 {
			time.Sleep(time.Millisecond * time.Duration(delta_timestamp-delta2))
		}
	}
}

// Closes flv file.
func (s *FlvStream) CloseFile() {
	s.FlvFile.Close()
}
