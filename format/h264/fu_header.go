package h264

import "github.com/searKing/rtp/codecs/h264"

const (
	FuHeaderMask      = 0xff
	FuHeaderOffset    = 0
	FuHeaderByteIndex = 1

	FuHeaderStartBitMask      = 1 << FuHeaderStartBitOffset
	FuHeaderStartBitOffset    = 7
	FuHeaderStartBitByteIndex = 0

	FuHeaderEndBitMask      = 1 << FuHeaderEndBitOffset
	FuHeaderEndBitOffset    = 6
	FuHeaderEndBitByteIndex = 0

	FuHeaderNalUnitTypeMask      = h264.NalUnitTypeMask
	FuHeaderNalUnitTypeOffset    = h264.NalUnitTypeOffset
	FuHeaderNalUnitTypeByteIndex = h264.NalUnitTypeByteIndex
)

// +---------------+
// |0|1|2|3|4|5|6|7|
// +-+-+-+-+-+-+-+-+
// |S|E|R|  Type   |
// +---------------+
// fuHeader
type FuHeader struct {
	StartBit bool
	EndBit   bool
	Type     h264.NalUnitType
}

func (h FuHeader) Byte() byte {
	b, _ := h.Marshal()
	return b[0]
}

func (h FuHeader) Marshal() ([]byte, error) {
	fuHeader := h.Type.Byte()
	if h.StartBit {
		fuHeader |= FuHeaderStartBitMask
	}

	if h.StartBit {
		fuHeader |= FuHeaderEndBitMask
	}

	return []byte{byte(fuHeader << FuHeaderOffset)}, nil
}

func (h *FuHeader) Unmarshal(buf []byte) error {
	h.StartBit = (buf[0] & FuHeaderStartBitMask) != 0
	h.EndBit = (buf[0] & FuHeaderEndBitMask) != 0
	_ = (&h.Type).Unmarshal(buf)
	return nil
}

func ParseFuHeader(rtpPayload []byte) FuHeader {
	var h FuHeader
	_ = (&h).Unmarshal(rtpPayload[FuHeaderByteIndex:])
	return h
}
