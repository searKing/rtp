package h264

import (
	"bytes"
	"fmt"
)

// H264Packet represents the H264 header that is stored in the payload of an RTP Packet
type Packet struct {
	// Required Header
	NalRefIdc   NalRefIdc
	NalUnitType NalUnitType

	Payload []byte
}

// Unmarshal parses the passed byte slice and stores the result in the VP8Packet this method is called upon
func (p *Packet) Unmarshal(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("invalid nil packet")
	}

	payloadLen := len(payload)

	if payloadLen < 1 {
		return fmt.Errorf("Payload is not large enough to container header")
	}

	payloadIndex := 0
	err := p.NalRefIdc.Unmarshal([]byte{payload[0]})
	if err != nil {
		return err
	}
	err = p.NalUnitType.Unmarshal([]byte{payload[0]})
	if err != nil {
		return err
	}

	payloadIndex++

	w := bytes.NewBuffer(nil)
	w.Write(payload[payloadIndex:])
	p.Payload = w.Bytes()
	return nil
}

// String helps with debugging by printing packet information in a readable way
func (p *Packet) String() string {
	out := "H264 Packet:\n"

	out += fmt.Sprintf("\tNalRefIdc: %v\n", p.NalRefIdc)
	out += fmt.Sprintf("\tNalUnitType: %s\n", p.NalUnitType)
	out += fmt.Sprintf("\tPayload Length: %d\n", len(p.Payload))

	return out
}

// Marshal serializes the header into bytes.
func (p *Packet) Marshal() ([]byte, error) {
	// avoid buf alloc
	w := bytes.NewBuffer(make([]byte, 0, p.MarshalSize()))
	byte0 := p.NalRefIdc.Byte()
	byte0 |= p.NalUnitType.Byte()
	w.WriteByte(byte0)
	w.Write(p.Payload)
	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (p *Packet) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	size := 1 + len(p.Payload)
	return size
}
