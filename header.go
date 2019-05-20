package rtp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Header represents an RTP packet header
type Header struct {
	Version        uint8
	Padding        bool
	Extension      bool
	Marker         bool
	PayloadType    uint8
	SequenceNumber uint16
	Timestamp      uint32
	SSRC           uint32
	CSRC           []uint32

	ExtensionProfile uint16
	ExtensionPayload []byte
}

// String helps with debugging by printing packet information in a readable way
func (h Header) String() string {
	out := "RTP Header:\n"

	out += fmt.Sprintf("\tVersion: %v\n", h.Version)
	out += fmt.Sprintf("\tPadding: %v\n", h.Padding)
	out += fmt.Sprintf("\tExtension: %v\n", h.Extension)
	out += fmt.Sprintf("\tMarker: %v\n", h.Marker)
	out += fmt.Sprintf("\tPayload Type: %d\n", h.PayloadType)
	out += fmt.Sprintf("\tSequence Number: %d\n", h.SequenceNumber)
	out += fmt.Sprintf("\tTimestamp: %d\n", h.Timestamp)
	out += fmt.Sprintf("\tSSRC: %d (%x)\n", h.SSRC, h.SSRC)
	out += fmt.Sprintf("\tCSRC Count: %d\n", len(h.CSRC))
	out += fmt.Sprintf("\tExtensionProfile: %d (%x)\n", h.ExtensionProfile, h.ExtensionProfile)
	out += fmt.Sprintf("\tExtensionPayload Length: %d\n", len(h.ExtensionPayload))

	return out
}

// Unmarshal parses the passed byte slice and stores the result in the Header this method is called upon
func (h *Header) Unmarshal(rawPacket []byte) error {
	if len(rawPacket) < headerLength {
		return fmt.Errorf("RTP header size insufficient; %d < %d", len(rawPacket), headerLength)
	}

	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|X|  CC   |M|     PT      |       sequence number         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                           timestamp                           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           synchronization source (SSRC) identifier            |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |            contributing source (CSRC) identifiers             |
	 * |                             ....                              |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	h.Version = rawPacket[0] >> versionShift & versionMask
	h.Padding = (rawPacket[0] >> paddingShift & paddingMask) > 0
	h.Extension = (rawPacket[0] >> extensionShift & extensionMask) > 0
	h.CSRC = make([]uint32, rawPacket[0]&ccMask)

	h.Marker = (rawPacket[1] >> markerShift & markerMask) > 0
	h.PayloadType = rawPacket[1] & ptMask

	h.SequenceNumber = binary.BigEndian.Uint16(rawPacket[seqNumOffset : seqNumOffset+seqNumLength])
	h.Timestamp = binary.BigEndian.Uint32(rawPacket[timestampOffset : timestampOffset+timestampLength])
	h.SSRC = binary.BigEndian.Uint32(rawPacket[ssrcOffset : ssrcOffset+ssrcLength])

	currOffset := csrcOffset + (len(h.CSRC) * csrcLength)
	if len(rawPacket) < currOffset {
		return fmt.Errorf("RTP header size insufficient; %d < %d", len(rawPacket), currOffset)
	}

	for i := range h.CSRC {
		offset := csrcOffset + (i * csrcLength)
		h.CSRC[i] = binary.BigEndian.Uint32(rawPacket[offset:])
	}

	if h.Extension {
		if len(rawPacket) < currOffset+4 {
			return fmt.Errorf("RTP header size insufficient for extension; %d < %d", len(rawPacket), currOffset)
		}

		h.ExtensionProfile = binary.BigEndian.Uint16(rawPacket[currOffset:])
		currOffset += 2
		extensionLength := int(binary.BigEndian.Uint16(rawPacket[currOffset:])) * 4
		currOffset += 2

		if len(rawPacket) < currOffset+extensionLength {
			return fmt.Errorf("RTP header size insufficient for extension length; %d < %d", len(rawPacket), currOffset+extensionLength)
		}

		h.ExtensionPayload = rawPacket[currOffset : currOffset+extensionLength]
		currOffset += len(h.ExtensionPayload)
	}

	return nil
}

// Marshal serializes the header into bytes.
func (h *Header) Marshal() ([]byte, error) {
	// avoid buf alloc
	buf := make([]byte, 0, h.MarshalSize())

	return h.marshal(buf)
}

// MarshalTo serializes the header and writes to the buffer.
func (h *Header) marshal(buf []byte) ([]byte, error) {
	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|X|  CC   |M|     PT      |       sequence number         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                           timestamp                           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           synchronization source (SSRC) identifier            |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |            contributing source (CSRC) identifiers             |
	 * |                             ....                              |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	size := h.MarshalSize()
	if size > cap(buf) {
		return nil, io.ErrShortBuffer
	}

	w := bytes.NewBuffer(buf)

	// The first byte contains the version, padding bit, extension bit, and csrc size
	byte0 := (h.Version << versionShift) | uint8(len(h.CSRC))
	if h.Padding {
		byte0 |= 1 << paddingShift
	}

	if h.Extension {
		byte0 |= 1 << extensionShift
	}

	w.WriteByte(byte0)

	// The second byte contains the marker bit and payload type.
	byte1 := h.PayloadType
	if h.Marker {
		byte1 |= 1 << markerShift
	}

	w.WriteByte(byte1)

	var word = make([]byte, 4)
	binary.BigEndian.PutUint16(word, h.SequenceNumber)
	w.Write(word[:2])
	binary.BigEndian.PutUint32(word, h.Timestamp)
	w.Write(word[:4])
	binary.BigEndian.PutUint32(word, h.SSRC)
	w.Write(word[:4])

	for _, csrc := range h.CSRC {
		binary.BigEndian.PutUint32(word, csrc)
		w.Write(word[:4])
	}

	if h.Extension {
		if len(h.ExtensionPayload)%4 != 0 {
			//the payload must be in 32-bit words.
			return w.Bytes(), io.ErrShortBuffer
		}
		extSize := uint16(len(h.ExtensionPayload) / 4)

		binary.BigEndian.PutUint16(word, h.ExtensionProfile)
		w.Write(word[:2])

		binary.BigEndian.PutUint16(word, extSize)
		w.Write(word[:2])

		w.Write(h.ExtensionPayload)
	}

	return w.Bytes(), nil
}

// MarshalSize returns the size of the header once marshaled.
func (h *Header) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	size := 12 + (len(h.CSRC) * csrcLength)

	if h.Extension {
		size += 4 + len(h.ExtensionPayload)
	}

	return size
}
