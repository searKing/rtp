package hevc

import "github.com/searKing/rtp/codecs/hevc"

/*
 * Table 1. – Summary of NAL unit types and the corresponding packet types in
 * rfc6184#section-5.2
 */

type RTPPacketType uint8

const (
	RTPPacketTypeMask      = hevc.NalUnitTypeMask
	RTPPacketTypeOffset    = hevc.NalUnitTypeOffset
	RTPPacketTypeByteIndex = hevc.NalUnitTypeByteIndex
)

const (
	RTPPacketTypeReserved     RTPPacketType = iota
	RTPPacketTypeNalUnitBegin               // 1
	_                                       // 2
	_                                       // 3
	_                                       // 4
	_                                       // 5
	_                                       // 6
	_                                       // 7
	_                                       // 8
	_                                       // 9
	_                                       // 10
	_                                       // 11
	_                                       // 12
	_                                       // 13
	_                                       // 14
	_                                       // 15
	_                                       // 16
	_                                       // 17
	_                                       // 18
	_                                       // 19
	_                                       // 20
	_                                       // 21
	_                                       // 22
	_                                       // 23
	_                                       // 24
	_                                       // 25
	_                                       // 26
	_                                       // 27
	_                                       // 28
	_                                       // 29
	_                                       // 30
	_                                       // 31
	_                                       // 32
	_                                       // 33
	_                                       // 34
	_                                       // 35
	_                                       // 36
	_                                       // 37
	_                                       // 38
	_                                       // 39
	_                                       // 40
	_                                       // 41
	_                                       // 42
	_                                       // 43
	_                                       // 44
	_                                       // 45
	_                                       // 46
	RTPPacketTypeNalUnitEnd                 // 47
	RTPPacketTypeAp                         // 48
	RTPPacketTypeFu                         // 49
	RTPPacketTypePaci                       // 50
	RTPPacketTypeReserved51                 // 51
	RTPPacketTypeReserved52                 // 52
	RTPPacketTypeReserved53                 // 53
	RTPPacketTypeReserved54                 // 54
	RTPPacketTypeReserved55                 // 55
	RTPPacketTypeReserved56                 // 56
	RTPPacketTypeReserved57                 // 57
	RTPPacketTypeReserved58                 // 58
	RTPPacketTypeReserved59                 // 59
	RTPPacketTypeReserved60                 // 60
	RTPPacketTypeReserved61                 // 61
	RTPPacketTypeReserved62                 // 62
	RTPPacketTypeReserved63                 // 63
)

func (t RTPPacketType) NalUnitType() hevc.NalUnitType {
	return hevc.NalUnitType(t)
}

func (t RTPPacketType) Byte() byte {
	b, _ := t.Marshal()
	return b[0]
}
func (t RTPPacketType) Marshal() ([]byte, error) {
	return []byte{byte(t << RTPPacketTypeOffset)}, nil
}

func (t *RTPPacketType) Unmarshal(buf []byte) error {
	*t = RTPPacketType((buf[0] & RTPPacketTypeMask) >> RTPPacketTypeOffset)
	return nil
}

func ParseRTPPacketType(payload []byte) RTPPacketType {
	var pt RTPPacketType
	_ = (&pt).Unmarshal(payload[RTPPacketTypeByteIndex:])
	return pt
}

func (t RTPPacketType) SingleNALUnitPacket() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|           PayloadHdr          |      DONL (conditional)       |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                                                               |
	//	|                  NAL unit payload data                        |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//		Figure 3: The Structure of a Single NAL Unit Packet
	if t >= RTPPacketTypeNalUnitBegin && t <= RTPPacketTypeNalUnitEnd {
		return true
	}
	return false
}
func (t RTPPacketType) AggregationPacket() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|    PayloadHdr (Type=48)       |                               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               |
	//	|                                                               |
	//	|             two or more aggregation units                     |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 4: The Structure of an Aggregation Packet

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//					:       DONL (conditional)      |   NALU size   |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|   NALU size   |                                               |
	//	+-+-+-+-+-+-+-+-+         NAL unit                              |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 5: The Structure of the First Aggregation Unit in an AP

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//					: DOND (cond)   |          NALU size            |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                                                               |
	//	|                       NAL unit                                |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	Figure 6: The Structure of an Aggregation Unit That Is Not the
	//	First Aggregation Unit in an AP
	switch t {
	case RTPPacketTypeAp:
		return true
	}
	return false
}
func (t RTPPacketType) FragmentationUnit() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|    PayloadHdr (Type=49)       |   FU header   | DONL (cond)   |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
	//	| DONL (cond)   |                                               |
	//	|-+-+-+-+-+-+-+-+                                               |
	//	|                         FU payload                            |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 9: The Structure of an FU

	//	+---------------+
	//	|0|1|2|3|4|5|6|7|
	//	+-+-+-+-+-+-+-+-+
	//	|S|E|  FuType   |
	//	+---------------+
	//
	//	Figure 10: The Structure of FU Header
	switch t {
	case RTPPacketTypeFu:
		return true
	}
	return false
}
func (t RTPPacketType) PACIPacket() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|    PayloadHdr (Type=50)       |A|   cType   | PHSsize |F0..2|Y|
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|        Payload Header Extension Structure (PHES)              |
	//	|=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=|
	//	|                                                               |
	//	|                  PACI payload: NAL unit                       |
	//	|                   . . .                                       |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 11: The Structure of a PACI

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|    PayloadHdr (Type=50)       |A|   cType   | PHSsize |F0..2|Y|
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|   TL0PICIDX   |   IrapPicID   |S|E|    RES    |               |
	//	|-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               |
	//	|                           ....                                |
	//	|               PACI payload: NAL unit                          |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 12: The Structure of a PACI with a PHES Containing a TSCI
	switch t {
	case RTPPacketTypePaci:
		return true
	}
	return false
}
func (t RTPPacketType) Reserved() bool {
	switch t {
	case RTPPacketTypeReserved,
		RTPPacketTypeReserved51,
		RTPPacketTypeReserved52,
		RTPPacketTypeReserved53,
		RTPPacketTypeReserved54,
		RTPPacketTypeReserved55,
		RTPPacketTypeReserved56,
		RTPPacketTypeReserved57,
		RTPPacketTypeReserved58,
		RTPPacketTypeReserved59,
		RTPPacketTypeReserved60,
		RTPPacketTypeReserved61,
		RTPPacketTypeReserved62,
		RTPPacketTypeReserved63:
		return true
	}
	return false
}

func (t RTPPacketType) PayloadHeaderSize() int {
	if t.SingleNALUnitPacket() {
		return 2
	}
	if t.AggregationPacket() {
		return 2
	}
	if t.FragmentationUnit() {
		return 2 + 1
	}
	if t.PACIPacket() {
		return 2 + 2 // + *PHSsize
	}
	return 2
}
