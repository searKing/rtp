package rtp

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

type PaddingTrailer struct {
	PaddingPayload []byte
}

// String helps with debugging by printing packet information in a readable way
func (t *PaddingTrailer) String() string {
	out := "RTP Padding Trailer:\n"

	out += fmt.Sprintf("\tPaddingPayload Length: %d\n", len(t.PaddingPayload))

	return out
}

// Unmarshal parses the passed byte slice and stores the result in the Header this method is called upon
func (t *PaddingTrailer) Unmarshal(padding []byte) error {
	if len(padding) < 1 {
		return nil
	}
	size := padding[len(padding)-1]

	// The last octet of the padding contains a count of how
	// many padding octets should be ignored, including itself.
	if len(padding) < int(size) {
		return fmt.Errorf("RTP padding trailer size insufficient; %d < %d", len(padding), size)
	}
	w := bytes.NewBuffer(nil)
	w.Write(padding[:size])
	t.PaddingPayload = w.Bytes()
	return nil
}

// Marshal serializes the header into bytes.
func (t *PaddingTrailer) Marshal() ([]byte, error) {
	buf := make([]byte, 0, t.MarshalSize())
	return t.marshal(buf)
}

// MarshalSize returns the size of the header once marshaled.
func (t *PaddingTrailer) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	if len(t.PaddingPayload) == 0 {
		return 0
	}
	// The last octet of the padding contains a count of how
	// many padding octets should be ignored, including itself.
	size := 1 + len(t.PaddingPayload)

	return size
}

// MarshalTo serializes the padding tailer and writes to the buffer.
func (t *PaddingTrailer) marshal(buf []byte) ([]byte, error) {
	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                             ....                              |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |              ....              |       padding length         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	size := t.MarshalSize()
	if size > math.MaxUint8 {
		return nil, fmt.Errorf("overflow padding payload, expect max %d, actual %d", math.MaxUint8, size)
	}
	if size == 0 {
		return nil, nil
	}

	if size > cap(buf) {
		return nil, io.ErrShortBuffer
	}

	w := bytes.NewBuffer(buf)
	w.Write(t.PaddingPayload)

	w.WriteByte(byte(size))

	return w.Bytes(), nil
}
