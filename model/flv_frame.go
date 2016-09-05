package model

import "github.com/zhangpeihao/goflv"

// Flv frame data
type FlvFrame struct {
	Header         *flv.TagHeader // Flv frame header.
	Frame          []byte         // Flv frame content.
	DeltaTimestamp uint32         // Flv frame delta constant.
}
