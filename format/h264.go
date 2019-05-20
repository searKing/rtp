package format

import "github.com/searKing/rtp/codecs/h264"

/*
 * Table 1. – Summary of NAL unit types and the corresponding packet types in
 * rfc6184#section-5.2
 */

type RTPPacketType uint8

const (
	RTPPacketTypeMask   = h264.NalUnitTypeMask
	RTPPacketTypeOffset = 0
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

func (t RTPPacketType) Byte() byte {
	b, _ := t.Marshal()
	return b[0]
}
func (t RTPPacketType) Marshal() ([]byte, error) {
	return []byte{byte(t << RTPPacketTypeOffset)}, nil
}

func (t *RTPPacketType) Unmarshal(payload []byte) error {
	*t = RTPPacketType(payload[0]&RTPPacketTypeMask) >> RTPPacketTypeOffset
	return nil
}

func ParseRTPPacketType(payload []byte) RTPPacketType {
	return RTPPacketType(payload[0] & RTPPacketTypeMask)
}

func (t RTPPacketType) SingleNALUnitPacket() bool {
	if t >= RTPPacketTypeNalUnitBegin && t <= RTPPacketTypeNalUnitEnd {
		return true
	}
	return false
}
func (t RTPPacketType) SingleTimeAggregationPacket() bool {
	switch t {
	case RTPPacketTypeStapA, RTPPacketTypeStapB:
		return true
	}
	return false
}
func (t RTPPacketType) MultiTimeAggregationPacket() bool {
	switch t {
	case RTPPacketTypeMtap16, RTPPacketTypeMtap24:
		return true
	}
	return false
}
func (t RTPPacketType) FragmentationUnit() bool {
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
