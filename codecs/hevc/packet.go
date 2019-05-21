package hevc

import (
	"bytes"
	"fmt"
)

// H264Packet represents the H264 header that is stored in the payload of an RTP Packet
type Packet struct {
	NalHeader NalHeader

	Payload []byte
}

// Unmarshal parses the passed byte slice and stores the result in the VP8Packet this method is called upon
func (p *Packet) Unmarshal(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("invalid nil packet")
	}

	payloadLen := len(payload)

	if payloadLen < 1 {
		return fmt.Errorf("payload is not large enough to container header")
	}

	payloadIndex := 0
	_ = (&p.NalHeader).Unmarshal(payload[NalHeaderByteIndex:])

	// jump to payload
	payloadIndex += 2

	w := bytes.NewBuffer(nil)
	w.Write(payload[payloadIndex:])
	p.Payload = w.Bytes()
	return nil
}

// Marshal serializes the header into bytes.
func (p *Packet) Marshal() ([]byte, error) {
	// avoid buf alloc
	w := bytes.NewBuffer(make([]byte, 0, p.MarshalSize()))
	w.Write(p.NalHeader.Bytes())

	w.Write(p.Payload)
	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (p *Packet) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	size := p.NalHeader.MarshalSize() + len(p.Payload)
	return size
}

// String helps with debugging by printing packet information in a readable way
func (p *Packet) String() string {
	out := "HEVC Packet:\n"

	out += fmt.Sprintf("\t%s\n", p.NalHeader)
	out += fmt.Sprintf("\tPayload Length: %d\n", len(p.Payload))

	return out
}
