package h264

import (
	"bytes"
	"fmt"
)

const (
	NalHeaderMask      = 0xff
	NalHeaderOffset    = 0
	NalHeaderByteIndex = 0
)

// NalHeader represents the nal header
type NalHeader struct {
	ForbiddenZeroBit ForbiddenZeroBit
	NalRefIdc        NalRefIdc
	NalUnitType      NalUnitType
}

// Unmarshal parses the passed byte slice and stores the result in the VP8NalHeader this method is called upon
func (h *NalHeader) Unmarshal(buf []byte) error {
	if buf == nil {
		return fmt.Errorf("invalid nil NalHeader")
	}

	bufSize := len(buf)

	if bufSize < 1 {
		return fmt.Errorf("buf is not large enough to container header")
	}
	_ = (&h.ForbiddenZeroBit).Unmarshal(buf[ForbiddenZeroBitByteIndex:])

	_ = (&h.NalRefIdc).Unmarshal(buf[NalRefIdcByteIndex:])

	_ = (&h.NalUnitType).Unmarshal(buf[NalUnitTypeByteIndex:])

	return nil
}

func (h NalHeader) Byte() byte {
	b, _ := h.Marshal()
	return b[0]
}

// Marshal serializes the header into bytes.
func (h NalHeader) Marshal() ([]byte, error) {
	// avoid buf alloc
	w := bytes.NewBuffer(make([]byte, 0, h.MarshalSize()))
	byte0 := h.ForbiddenZeroBit.Byte()
	byte0 |= h.NalRefIdc.Byte()
	byte0 |= h.NalUnitType.Byte()
	w.WriteByte(byte0)
	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (h NalHeader) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	return 1
}

// String helps with debugging by printing NalHeader information in a readable way
func (h NalHeader) String() string {
	out := "H264 NalHeader:\n"

	out += fmt.Sprintf("\tForbiddenZeroBit: %v\n", h.ForbiddenZeroBit)
	out += fmt.Sprintf("\tNalRefIdc: %v\n", h.NalRefIdc)
	out += fmt.Sprintf("\tNalUnitType: %s\n", h.NalUnitType)

	return out
}

func ParseNalHeader(nalu []byte) NalHeader {
	var h NalHeader
	_ = (&h).Unmarshal(nalu[NalHeaderByteIndex:])
	return h
}
