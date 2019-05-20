package rtp

import (
	"bytes"
	"fmt"
	"io"
)

// Packet represents an RTP Packet
// NOTE: Raw is populated by Marshal/Unmarshal and should not be modified
type Packet struct {
	Header         Header
	Payload        []byte
	PaddingTrailer PaddingTrailer
}

const (
	headerLength    = 4
	versionShift    = 6
	versionMask     = 0x3
	paddingShift    = 5
	paddingMask     = 0x1
	extensionShift  = 4
	extensionMask   = 0x1
	ccMask          = 0xF
	markerShift     = 7
	markerMask      = 0x1
	ptMask          = 0x7F
	seqNumOffset    = 2
	seqNumLength    = 2
	timestampOffset = 4
	timestampLength = 4
	ssrcOffset      = 8
	ssrcLength      = 4
	csrcOffset      = 12
	csrcLength      = 4
)

// String helps with debugging by printing packet information in a readable way
func (p Packet) String() string {
	out := "RTP PACKET:\n"

	out += fmt.Sprintf("%s\n", p.Header)
	out += fmt.Sprintf("\tPayload Length: %d\n", len(p.Payload))
	out += fmt.Sprintf("%s", p.PaddingTrailer)

	return out
}

// Unmarshal parses the passed byte slice and stores the result in the Packet this method is called upon
func (p *Packet) Unmarshal(rawPacket []byte) error {
	if err := p.Header.Unmarshal(rawPacket); err != nil {
		return err
	}
	if p.Header.Padding {
		if err := p.PaddingTrailer.Unmarshal(rawPacket); err != nil {
			return err
		}
	}
	payloadWithPadding := rawPacket[p.Header.MarshalSize():]
	p.Payload = payloadWithPadding[:len(payloadWithPadding)-p.PaddingTrailer.MarshalSize()]
	return nil
}

// Marshal serializes the packet into bytes.
func (p *Packet) Marshal() (buf []byte, err error) {
	buf = make([]byte, 0, p.MarshalSize())

	return p.marshal(buf)
}

// MarshalTo serializes the packet and writes to the buffer.
func (p *Packet) marshal(buf []byte) ([]byte, error) {
	// Make sure the buffer is large enough to hold the packet.

	if cap(buf) < p.MarshalSize() {
		return nil, io.ErrShortBuffer
	}

	h, err := p.Header.Marshal()
	if err != nil {
		return nil, err
	}
	w := bytes.NewBuffer(buf)
	w.Write(h)

	w.Write(p.Payload)
	if p.Header.Padding {
		t, err := p.PaddingTrailer.Marshal()
		if err != nil {
			return w.Bytes(), err
		}
		w.Write(t)
	}
	return w.Bytes(), nil
}

// MarshalSize returns the size of the packet once marshaled.
func (p *Packet) MarshalSize() int {
	return p.Header.MarshalSize() + len(p.Payload) + p.PaddingTrailer.MarshalSize()
}
