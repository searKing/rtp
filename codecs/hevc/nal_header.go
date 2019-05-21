package hevc

import (
	"bytes"
	"fmt"
	"github.com/searKing/rtp/codecs/h264"
)

const (
	NalHeaderMask      = 0xff
	NalHeaderOffset    = 0
	NalHeaderByteIndex = 0
)

//	 0                   1
//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|F|   Type    |  LayerId  | TID |
//	+-------------+-----------------+

// NalHeader represents the HEVC header
type NalHeader struct {
	ForbiddenZeroBit h264.ForbiddenZeroBit
	NalUnitType      NalUnitType
	NalLayerId       NalLayerId
	NalTemporalId    NalTemporalId
}

// Unmarshal parses the passed byte slice and stores the result in the VP8NalHeader this method is called upon
func (h *NalHeader) Unmarshal(buf []byte) error {
	if buf == nil {
		return fmt.Errorf("invalid nil packet")
	}

	payloadLen := len(buf)

	if payloadLen < 1 {
		return fmt.Errorf("payload is not large enough to container header")
	}

	_ = (&h.ForbiddenZeroBit).Unmarshal(buf[h264.ForbiddenZeroBitByteIndex:])

	_ = (&h.NalUnitType).Unmarshal(buf[NalUnitTypeByteIndex:])

	_ = (&h.NalLayerId).Unmarshal(buf[NalLayerIdByteIndex:])
	_ = (&h.NalTemporalId).Unmarshal(buf[NalTemporalIdByteIndex:])

	return nil
}

func (h NalHeader) Bytes() []byte {
	b, _ := h.Marshal()
	return b
}

// Marshal serializes the header into bytes.
func (h NalHeader) Marshal() ([]byte, error) {
	// avoid buf alloc
	w := bytes.NewBuffer(make([]byte, 0, h.MarshalSize()))
	layerIds, _ := h.NalLayerId.Marshal()
	byte0 := h.ForbiddenZeroBit.Byte()
	byte0 |= h.NalUnitType.Byte()
	byte0 |= layerIds[0]

	byte1 := h.NalTemporalId.Byte() | layerIds[1]

	w.WriteByte(byte0)
	w.WriteByte(byte1)
	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (h NalHeader) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	return 2
}

// String helps with debugging by printing packet information in a readable way
func (h NalHeader) String() string {
	out := "HEVC NalHeader:\n"

	out += fmt.Sprintf("\tForbiddenZeroBit: %v\n", h.ForbiddenZeroBit)
	out += fmt.Sprintf("\tNalUnitType: %s\n", h.NalUnitType)
	out += fmt.Sprintf("\tNalLayerId: %d\n", h.NalLayerId)
	out += fmt.Sprintf("\tNalTemporalId: %d\n", h.NalTemporalId)

	return out
}

func ParseNalHeader(nalu []byte) NalHeader {
	var h NalHeader
	_ = (&h).Unmarshal(nalu[NalHeaderByteIndex:])
	return h
}
