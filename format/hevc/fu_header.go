package hevc

import "github.com/searKing/rtp/codecs/hevc"

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

	FuHeaderFuTypeMask      = hevc.NalUnitTypeMask
	FuHeaderFuTypeOffset    = hevc.NalUnitTypeOffset
	FuHeaderFuTypeByteIndex = hevc.NalUnitTypeByteIndex
)

// +---------------+
// |0|1|2|3|4|5|6|7|
// +-+-+-+-+-+-+-+-+
// |S|E|  FuType   |
// +---------------+
// fuHeader
type FuHeader struct {
	StartBit bool
	EndBit   bool
	FuType   hevc.NalUnitType
}

func (h FuHeader) Byte() byte {
	b, _ := h.Marshal()
	return b[0]
}

func (h FuHeader) Marshal() ([]byte, error) {
	fuHeader := h.FuType.Byte()
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
	_ = (&h.FuType).Unmarshal(buf)
	return nil
}

func ParseFuHeader(rtpPayload []byte) FuHeader {
	var h FuHeader
	_ = (&h).Unmarshal(rtpPayload[FuHeaderByteIndex:])
	return h
}
