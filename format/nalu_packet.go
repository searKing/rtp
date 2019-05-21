package format

import (
	"bytes"
	"encoding/binary"
	h264_codec "github.com/searKing/rtp/codecs/h264"
	hevc_codec "github.com/searKing/rtp/codecs/hevc"
	"github.com/searKing/rtp/format/h264"
	"github.com/searKing/rtp/format/hevc"
	"math"
)

func naluPacket(maxPayloadSize int, nals [][]byte, skipAggregate bool, h264NotHevc bool) [][]byte {
	var packetedNals [][]byte
	payloadHeaderSize := func() int {
		if h264NotHevc {
			return h264.RTPPacketTypeStapA.HeaderSize()
		}
		return hevc.RTPPacketTypeAp.PayloadHeaderSize()
	}()

	var nalbuffersSize int
	var nalbuffers [][]byte

	flushBufferedNals := func() {
		defer func() {
			nalbuffers = nil
			nalbuffersSize = 0
		}()
		// Flush buffered NAL units if the current unit doesn't fit
		if len(nalbuffers) == 0 {
			return
		}
		// Fragment
		if len(nalbuffers) == 1 {
			packetedNals = append(packetedNals, tryFragmentNaluIfNecessary(maxPayloadSize, nalbuffers[0], h264NotHevc)...)
			return
		}

		// Aggregate
		packetedNals = append(packetedNals, tryAggregateNalus(nalbuffers, h264NotHevc)...)
	}

	for _, nal := range nals {
		if len(nal) == 0 || len(nal) > math.MaxUint16 {
			// invalid nal
			continue
		}

		if len(nal) <= maxPayloadSize {
			// If the NAL unit fits including the
			// framing (2 bytes length, plus 1/2 bytes for the STAP-A/AP marker),
			// write the unit to the buffer as a STAP-A/AP packet, otherwise flush
			// and send as single NAL.

			//	 0                   1                   2                   3
			//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
			//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//	|F|NRI|  Type   |        NAL unit size          |
			//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

			//	 0                   1                   2                   3
			//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
			//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//	|   PayloadHdr (Type=48)        |          NALU 1 Size          |
			//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			//	|            NALU 1 HDR         |
			//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			payloadSize := payloadHeaderSize + 2 + len(nal)
			if nalbuffersSize+payloadSize <= maxPayloadSize && !skipAggregate {
				// aggregate
				nalbuffers = append(nalbuffers, nal)
				nalbuffersSize += len(nal)
				continue
			}
		}

		flushBufferedNals()
		nalbuffers = append(nalbuffers, nal)
		nalbuffersSize += len(nal)
		// fragment this nalu
		if len(nal) > maxPayloadSize {
			flushBufferedNals()
		}
	}
	flushBufferedNals()

	return packetedNals

}

func tryFragmentNaluIfNecessary(maxPayloadSize int, nalu []byte, h264NotHevc bool) [][]byte {
	var fragmentedNals [][]byte
	headerSize := func() int {
		if h264NotHevc {
			return h264.RTPPacketTypeFuA.HeaderSize()
		}
		return hevc.RTPPacketTypeFu.PayloadHeaderSize()
	}()

	// Single NALU
	if len(nalu) <= maxPayloadSize {
		out := make([]byte, len(nalu))
		copy(out, nalu)
		fragmentedNals = append(fragmentedNals, out)
		return fragmentedNals
	}

	// FU-A
	maxFragmentSize := maxPayloadSize - headerSize // nalu without nal header(1Byte)

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
		return fragmentedNals
	}
	currentNalDataFragmentSize := min(maxFragmentSize, naluDataRemaining)

	var h264FuIndicator h264.FuIndicator
	var h264FuHeader h264.FuHeader

	var hevcPayloadHdr hevc.PayloadHdr
	var hevcFuHeader hevc.FuHeader
	if h264NotHevc {
		h264FuIndicator, h264FuHeader = initNaluH264Fu(nalu)
	} else {
		hevcPayloadHdr, hevcFuHeader = initNaluHEVCFu(nalu)
	}

	for naluDataRemaining > 0 {
		if naluDataRemaining == naluDataLength {
			// Set start bit
			h264FuHeader.StartBit = true
			h264FuHeader.EndBit = false
			hevcFuHeader.StartBit = true
			hevcFuHeader.EndBit = false
		} else if naluDataRemaining-currentNalDataFragmentSize == 0 {
			// Set end bit
			h264FuHeader.StartBit = false
			h264FuHeader.EndBit = true
			hevcFuHeader.StartBit = false
			hevcFuHeader.EndBit = true
		} else {
			h264FuHeader.StartBit = false
			h264FuHeader.EndBit = false
			hevcFuHeader.StartBit = false
			hevcFuHeader.EndBit = false
		}

		w := bytes.NewBuffer(make([]byte, headerSize+currentNalDataFragmentSize))
		if h264NotHevc {
			w.WriteByte(h264FuIndicator.Byte())
			w.WriteByte(h264FuHeader.Byte())
		} else {
			w.Write(hevcPayloadHdr.Bytes())
			w.WriteByte(hevcFuHeader.Byte())
		}
		w.Write(naluData[naluDataIndex : naluDataIndex+currentNalDataFragmentSize])
		fragmentedNals = append(fragmentedNals, w.Bytes())

		naluDataRemaining -= currentNalDataFragmentSize
		naluDataIndex += currentNalDataFragmentSize
	}
	return fragmentedNals
}
func tryAggregateNalus(nalbuffers [][]byte, h264NotHevc bool) [][]byte {
	var aggregatedNals [][]byte

	for _, nal := range nalbuffers {
		//aggregate buffered nalus
		w := bytes.NewBuffer(nil)
		if h264NotHevc {
			fuIndicator := h264.FuIndicator{}
			fuIndicator.NalUnitType = h264.RTPPacketTypeStapA.NalUnitType()
			w.WriteByte(fuIndicator.Byte())
		} else {
			payloadHdr := hevc.PayloadHdr{}
			payloadHdr.NalUnitType = hevc.RTPPacketTypeAp.NalUnitType()
			payloadHdr.NalTemporalId = 1
			w.Write(payloadHdr.Bytes())
		}
		var word = make([]byte, 4)
		binary.BigEndian.PutUint16(word, uint16(len(nal)))
		w.Write(word[:2])
		w.Write(nal)
		aggregatedNals = append(aggregatedNals, w.Bytes())
	}
	return aggregatedNals
}
func initNaluH264Fu(nalu []byte) (h264.FuIndicator, h264.FuHeader) {
	naluHeader := h264_codec.ParseNalHeader(nalu)
	// +---------------+
	// |0|1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+-+
	// |F|NRI|  Type   |
	// +---------------+
	// fuIndicator
	fuIndicator := h264.FuIndicator{
		NalHeader: naluHeader,
	}
	fuIndicator.NalUnitType = h264.RTPPacketTypeFuA.NalUnitType()

	// +---------------+
	// |0|1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+-+
	// |S|E|R|  Type   |
	// +---------------+
	// fuHeader
	fuHeader := h264.FuHeader{}
	fuHeader.Type = naluHeader.NalUnitType
	return fuIndicator, fuHeader
}

func initNaluHEVCFu(nalu []byte) (hevc.PayloadHdr, hevc.FuHeader) {
	naluHeader := hevc_codec.ParseNalHeader(nalu)
	payloadHdr := hevc.PayloadHdr{}
	payloadHdr.NalUnitType = hevc.RTPPacketTypeFu.NalUnitType()
	payloadHdr.NalTemporalId = 1
	// +---------------+
	// |0|1|2|3|4|5|6|7|
	// +-+-+-+-+-+-+-+-+
	// |S|E|R|  Type   |
	// +---------------+
	// fuHeader
	fuHeader := hevc.FuHeader{
		FuType: naluHeader.NalUnitType,
	}
	return payloadHdr, fuHeader
}
