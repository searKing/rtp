package hevc

import "github.com/searKing/rtp/codecs/hevc"

const (
	PayloadHdrMask      = 0xffff
	PayloadHdrOffset    = 0
	PayloadHdrByteIndex = 0
)

type PayloadHdr struct {
	hevc.NalHeader
}

func (hdr PayloadHdr) Bytes() []byte {
	b, _ := hdr.Marshal()
	return b
}

func (hdr PayloadHdr) Marshal() ([]byte, error) {
	nalHeader := hdr.NalHeader.Bytes()
	return nalHeader, nil
}

func (hdr *PayloadHdr) Unmarshal(buf []byte) error {
	_ = (&hdr.NalHeader).Unmarshal(buf)
	return nil
}

func ParsePayloadHdr(rtpPayload []byte) PayloadHdr {
	var h PayloadHdr
	_ = (&h).Unmarshal(rtpPayload[PayloadHdrByteIndex:])
	return h
}
