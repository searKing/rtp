package h264

import "github.com/searKing/rtp/codecs/h264"

/*
 * Table 1. – Summary of NAL unit types and the corresponding packet types in
 * rfc6184#section-5.2
 */

type RTPPacketType uint8

const (
	RTPPacketTypeMask      = h264.NalUnitTypeMask
	RTPPacketTypeOffset    = h264.NalUnitTypeOffset
	RTPPacketTypeByteIndex = h264.NalUnitTypeByteIndex
)

// 	Table 1.  Summary of NAL unit types and the corresponding packet types
//
// 	NAL Unit  Packet    Packet Type Name               Section
// 	Type      Type
// 	-------------------------------------------------------------
//	0        reserved                                     -
//	1-23     NAL unit  Single NAL unit packet             5.6
// 	24       STAP-A    Single-time aggregation packet     5.7.1
//	25       STAP-B    Single-time aggregation packet     5.7.1
//	26       MTAP16    Multi-time aggregation packet      5.7.2
//	27       MTAP24    Multi-time aggregation packet      5.7.2
//	28       FU-A      Fragmentation unit                 5.8
//	29       FU-B      Fragmentation unit                 5.8
//	30-31    reserved                                     -
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
	RTPPacketTypeNalUnitEnd                 // 23
	RTPPacketTypeStapA                      // 24
	RTPPacketTypeStapB                      // 25
	RTPPacketTypeMtap16                     // 26
	RTPPacketTypeMtap24                     // 27
	RTPPacketTypeFuA                        // 28
	RTPPacketTypeFuB                        // 29
	RTPPacketTypeReserved30                 //30
	RTPPacketTypeReserved31                 //31
)

func (t RTPPacketType) NalUnitType() h264.NalUnitType {
	return h264.NalUnitType(t)
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
	//	|F|NRI|  Type   |                                               |
	//	+-+-+-+-+-+-+-+-+                                               |
	//	|                                                               |
	//	|               Bytes 2..n of a single NAL unit                 |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//		Figure 2.  RTP payload format for single NAL unit packet
	if t >= RTPPacketTypeNalUnitBegin && t <= RTPPacketTypeNalUnitEnd {
		return true
	}
	return false
}
func (t RTPPacketType) SingleTimeAggregationPacket() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|F|NRI|  Type   |                                               |
	//	+-+-+-+-+-+-+-+-+                                               |
	//	|                                                               |
	//	|             one or more aggregation units                     |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//		Figure 3.  RTP payload format for aggregation packets

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|F|NRI|  Type   |                                               |
	//	+-+-+-+-+-+-+-+-+                                               |
	//	|                                                               |
	//	|                single-time aggregation units                  |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 4.  Payload format for STAP-A

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|F|NRI|  Type   |  decoding order number (DON)  |               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               |
	//	|                                                               |
	//	|                single-time aggregation units                  |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 5.  Payload format for STAP-B

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|F|NRI|  Type   |        NAL unit size          |               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               |
	//	|                                                               |
	//	|                           NAL unit                            |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 6.  Structure for single-time aggregation unit
	switch t {
	case RTPPacketTypeStapA, RTPPacketTypeStapB:
		return true
	}
	return false
}
func (t RTPPacketType) MultiTimeAggregationPacket() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|MTAP16 NAL HDR |:  decoding order number base   |               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+               |
	//	|                                                               |
	//	|                 multi-time aggregation units                  |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 9.  NAL unit payload format for MTAPs

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	:        NAL unit size          |      DOND     |  TS offset    |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|  TS offset    |                                               |
	//	+-+-+-+-+-+-+-+-+              NAL unit                         |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 10.  Multi-time aggregation unit for MTAP16

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	:        NAL unit size         |      DOND     |  TS offset    |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|         TS offset             |                               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               |
	//	|                              NAL unit                         |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//	Figure 11.  Multi-time aggregation unit for MTAP24
	switch t {
	case RTPPacketTypeMtap16, RTPPacketTypeMtap24:
		return true
	}
	return false
}
func (t RTPPacketType) FragmentationUnit() bool {
	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	| FU indicator  |   FU header   |                               |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                               |
	//	|                                                               |
	//	|                         FU payload                            |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//		Figure 14.  RTP payload format for FU-A

	//	 0                   1                   2                   3
	//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	| FU indicator  |   FU header   |               DON             |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
	//	|                                                               |
	//	|                         FU payload                            |
	//	|                                                               |
	//	|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//	|                               :...OPTIONAL RTP padding        |
	//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//
	//		Figure 15.  RTP payload format for FU-B
	//	The FU indicator octet has the following format:
	//
	//	+---------------+
	//	|0|1|2|3|4|5|6|7|
	//	+-+-+-+-+-+-+-+-+
	//	|F|NRI|  Type   |
	//	+---------------+

	//	The FU header has the following format:
	//
	//	+---------------+
	//	|0|1|2|3|4|5|6|7|
	//	+-+-+-+-+-+-+-+-+
	//	|S|E|R|  Type   |
	//	+---------------+
	switch t {
	case RTPPacketTypeFuA, RTPPacketTypeFuB:
		return true
	}
	return false
}
func (t RTPPacketType) Reserved() bool {
	switch t {
	case RTPPacketTypeReserved, RTPPacketTypeReserved30, RTPPacketTypeReserved31:
		return true
	}
	return false
}

func (t RTPPacketType) HeaderSize() int {
	if t.SingleNALUnitPacket() {
		return 1
	}
	switch t {
	case RTPPacketTypeStapA:
		return 1
	case RTPPacketTypeStapB:
		return 1 + 2
	case RTPPacketTypeMtap16, RTPPacketTypeMtap24:
		return 1 + 2
	case RTPPacketTypeFuA:
		return 1 + 1
	case RTPPacketTypeFuB:
		return 1 + 1 + 2
	default:
		return 1
	}
}
