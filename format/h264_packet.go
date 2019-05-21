package format

// H264Payloader payloads H264 packets
type H264Payloader struct {
	SkipAggregate bool
}

const (
	fuaHeaderSize = 2

	fuHeaderStartBitMask   = 1 << fuHeaderStartBitOffset
	fuHeaderStartBitOffset = 7
	fuHeaderEndBitMask     = 1 << fuHeaderEndBitOffset
	fuHeaderEndBitOffset   = 6

	fuHeaderNalUnitTypeMask   = 0x1f
	fuHeaderNalUnitTypeOffset = 0
)

// Payload fragments a H264 packet across one or more byte arrays
// ffmpeg/libavformat/rtpenc_h264_hevc.c nal_send
func (p *H264Payloader) Payload(maxPayloadSize int, payload []byte) [][]byte {

	if payload == nil {
		return nil
	}
	var nalubuffer [][]byte
	emitNalus(payload, func(nalu []byte) {
		nalubuffer = append(nalubuffer, nalu)
	})
	return naluPacket(maxPayloadSize, nalubuffer, p.SkipAggregate, false)
}

// traversal nals and emit when a nalu is meet
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
