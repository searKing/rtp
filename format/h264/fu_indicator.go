package h264

import "github.com/searKing/rtp/codecs/h264"

const (
	FuIndicatorMask      = h264.NalHeaderMask
	FuIndicatorOffset    = h264.NalHeaderOffset
	FuIndicatorByteIndex = h264.NalHeaderByteIndex
)

// +---------------+
// |0|1|2|3|4|5|6|7|
// +-+-+-+-+-+-+-+-+
// |F|NRI|  Type   |
// +---------------+
// FuIndicator
type FuIndicator struct {
	h264.NalHeader
}

func (h FuIndicator) Byte() byte {
	return h.NalHeader.Byte()
}

func (h FuIndicator) Marshal() ([]byte, error) {
	return h.NalHeader.Marshal()
}

// MarshalSize returns the size of the header once marshaled.
func (h FuIndicator) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	return h.NalHeader.MarshalSize()
}

func (h *FuIndicator) Unmarshal(buf []byte) error {
	_ = (&h.NalHeader).Unmarshal(buf)
	return nil
}

func ParseFuIndicator(rtpPayload []byte) FuIndicator {
	var h FuIndicator
	_ = (&h).Unmarshal(rtpPayload[FuIndicatorByteIndex:])
	return h
}
