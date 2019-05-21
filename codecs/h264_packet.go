package codecs

import (
	"bytes"
	"fmt"
	"github.com/searKing/rtp/codecs/h264"
	"github.com/searKing/rtp/format"
)

// H264Payloader payloads H264 packets
type H264Payloader struct{}

const (
	fuaHeaderSize = 2

	fuHeaderStartBitMask   = 1 << fuHeaderStartBitOffset
	fuHeaderStartBitOffset = 7
	fuHeaderEndBitMask     = 1 << fuHeaderEndBitOffset
	fuHeaderEndBitOffset   = 6

	fuHeaderNalUnitTypeMask   = 0x1f
	fuHeaderNalUnitTypeOffset = 0
)

func emitNalus(nals []byte, emit func([]byte)) {
	// for leading code, such as 0x 00 00 00 01, 0x00's count must be more than one, 0x00 00 01 at least
	nextInd := func(nalu []byte, start int) (indStart int, indLen int) {
		zeroCount := 0

		// 0x00 00 00 01
		for i, b := range nalu[start:] {
			if b == 0 {
				zeroCount++
				continue
			} else if b == 1 {
				if zeroCount >= 2 {
					return start + i - zeroCount, zeroCount + 1
				}
			}
			zeroCount = 0
		}
		return -1, -1
	}
	nextIndStart := 0
	nextIndLen := 0
	for {
		prevStart := nextIndStart + nextIndLen
		nextIndStart, nextIndLen = nextInd(nals, prevStart)
		if nextIndStart == -1 {
			// Emit until end of stream, no end indicator found
			emit(nals[prevStart:])
			return
		}
		if prevStart == 0 {
			continue
		}
		emit(nals[prevStart:nextIndStart])
	}
}

// Payload fragments a H264 packet across one or more byte arrays
func (p *H264Payloader) Payload(mtu int, payload []byte) [][]byte {

	var payloads [][]byte
	if payload == nil {
		return payloads
	}

	emitNalus(payload, func(nalu []byte) {
		//naluType := nalu[0] & h264.NalUnitTypeMask
		naluType := h264.ParseNalUnitType(nalu)
		naluRefIdc := h264.ParseNalRefIdc(nalu)
		if naluType == h264.NalUnitTypeAud || naluType == h264.NalUnitTypeFillerData {
			return
		}

		// Single NALU
		if len(nalu) <= mtu {
			out := make([]byte, len(nalu))
			copy(out, nalu)
			payloads = append(payloads, out)
			return
		}

		// FU-A
		maxFragmentSize := mtu - fuaHeaderSize

		// The FU payload consists of fragments of the payload of the fragmented
		// NAL unit so that if the fragmentation unit payloads of consecutive
		// FUs are sequentially concatenated, the payload of the fragmented NAL
		// unit can be reconstructed.  The NAL unit type octet of the fragmented
		// NAL unit is not included as such in the fragmentation unit payload,
		// 	but rather the information of the NAL unit type octet of the
		// fragmented NAL unit is conveyed in the F and NRI fields of the FU
		// indicator octet of the fragmentation unit and in the type field of
		// the FU header.  An FU payload MAY have any number of octets and MAY
		// be empty.

		naluData := nalu
		// According to the RFC, the first octet is skipped due to redundant information
		naluDataIndex := 1
		naluDataLength := len(nalu) - naluDataIndex
		naluDataRemaining := naluDataLength

		if min(maxFragmentSize, naluDataRemaining) <= 0 {
			return
		}

		for naluDataRemaining > 0 {
			currentFragmentSize := min(maxFragmentSize, naluDataRemaining)

			// +---------------+
			// |0|1|2|3|4|5|6|7|
			// +-+-+-+-+-+-+-+-+
			// |F|NRI|  Type   |
			// +---------------+
			// fuIndicator
			fuIndicator := format.RTPPacketTypeFuA.Byte()
			fuIndicator |= naluRefIdc.Byte()

			// +---------------+
			// |0|1|2|3|4|5|6|7|
			// +-+-+-+-+-+-+-+-+
			// |S|E|R|  Type   |
			// +---------------+
			// fuHeader
			fuHeader := naluType.Byte()

			if naluDataRemaining == naluDataLength {
				// Set start bit
				fuHeader |= fuHeaderStartBitMask
			} else if naluDataRemaining-currentFragmentSize == 0 {
				// Set end bit
				fuHeader |= fuHeaderEndBitMask
			}
			w := bytes.NewBuffer(make([]byte, fuaHeaderSize+currentFragmentSize))
			w.WriteByte(fuIndicator)
			w.WriteByte(fuHeader)
			w.Write(naluData[naluDataIndex : naluDataIndex+currentFragmentSize])
			payloads = append(payloads, w.Bytes())

			naluDataRemaining -= currentFragmentSize
			naluDataIndex += currentFragmentSize
		}

	})

	return payloads
}

// H264Packet represents the H264 header that is stored in the payload of an RTP Packet
type H264Packet struct {
	// Required Header
	NalRefIdc   h264.NalRefIdc
	NalUnitType h264.NalUnitType

	Payload []byte
}

// Unmarshal parses the passed byte slice and stores the result in the VP8Packet this method is called upon
func (p *H264Packet) Unmarshal(payload []byte) error {
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
func (p *H264Packet) String() string {
	out := "H264 Packet:\n"

	out += fmt.Sprintf("\tNalRefIdc: %v\n", p.NalRefIdc)
	out += fmt.Sprintf("\tNalUnitType: %s\n", p.NalUnitType)
	out += fmt.Sprintf("\tPayload Length: %d\n", len(p.Payload))

	return out
}

// Marshal serializes the header into bytes.
func (p *H264Packet) Marshal() ([]byte, error) {
	// avoid buf alloc
	w := bytes.NewBuffer(make([]byte, 0, p.MarshalSize()))
	byte0 := p.NalRefIdc.Byte()
	byte0 |= p.NalUnitType.Byte()
	w.WriteByte(byte0)
	w.Write(p.Payload)
	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (p *H264Packet) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	size := 1 + len(p.Payload)
	return size
}
