package h264

import "fmt"

//	+---------------+
//	|0|1|2|3|4|5|6|7|
//	+-+-+-+-+-+-+-+-+
//	|F|NRI|  Type   |
//	+---------------+

const (
	NalUnitTypeMask      = 0x1f
	NalUnitTypeOffset    = 0
	NalUnitTypeByteIndex = 0
)

// ffmpeg/libavcodec/h264.h
/*
 * Table 7-1 – NAL unit type codes, syntax element categories, and NAL unit type classes in
 * T-REC-H.264-201704
 */
type NalUnitType uint8

const (
	NalUnitTypeUnspecified     NalUnitType = iota // 0
	NalUnitTypeSlice                              // 1
	NalUnitTypeDpa                                // 2
	NalUnitTypeDpb                                // 3
	NalUnitTypeDpc                                // 4
	NalUnitTypeIdrSlice                           // 5
	NalUnitTypeSei                                // 6
	NalUnitTypeSps                                // 7
	NalUnitTypePps                                // 8
	NalUnitTypeAud                                // 9
	NalUnitTypeEndSequence                        // 10
	NalUnitTypeEndStream                          // 11
	NalUnitTypeFillerData                         // 12
	NalUnitTypeSpsExt                             // 13
	NalUnitTypePrefix                             // 14
	NalUnitTypeSubSps                             // 15
	NalUnitTypeDps                                // 16
	NalUnitTypeReserved17                         // 17
	NalUnitTypeReserved18                         // 18
	NalUnitTypeAuxiliarySlice                     // 19
	NalUnitTypeExtenSlice                         // 20
	NalUnitTypeDepthExtenSlice                    // 21
	NalUnitTypeReserved22                         // 22
	NalUnitTypeReserved23                         // 23
	NalUnitTypeUnspecified24                      // 24
	NalUnitTypeUnspecified25                      // 25
	NalUnitTypeUnspecified26                      // 26
	NalUnitTypeUnspecified27                      // 27
	NalUnitTypeUnspecified28                      // 28
	NalUnitTypeUnspecified29                      // 29
	NalUnitTypeUnspecified30                      // 30
	NalUnitTypeUnspecified31                      // 31
)

func (naluType NalUnitType) Byte() byte {
	b, _ := naluType.Marshal()
	return b[0]
}

func (naluType NalUnitType) Marshal() ([]byte, error) {
	return []byte{byte(naluType << NalUnitTypeOffset)}, nil
}

func (naluType *NalUnitType) Unmarshal(buf []byte) error {
	*naluType = NalUnitType((buf[0] & NalUnitTypeMask) >> NalUnitTypeOffset)
	return nil
}

func (naluType NalUnitType) String() string {
	switch naluType {
	case NalUnitTypeUnspecified:
		return "unspecified"
	case NalUnitTypeSlice:
		return "slice"
	case NalUnitTypeDpa:
		return "dpa"
	case NalUnitTypeDpb:
		return "dpb"
	case NalUnitTypeDpc:
		return "dpc"
	case NalUnitTypeIdrSlice:
		return "idr slice"
	case NalUnitTypeSei:
		return "sei"
	case NalUnitTypeSps:
		return "sps"
	case NalUnitTypePps:
		return "pps"
	case NalUnitTypeAud:
		return "aud"
	case NalUnitTypeEndSequence:
		return "end sequence"
	case NalUnitTypeEndStream:
		return "end stream"
	case NalUnitTypeFillerData:
		return "filler data"
	case NalUnitTypeSpsExt:
		return "sps ext"
	case NalUnitTypePrefix:
		return "prefix"
	case NalUnitTypeSubSps:
		return "sub sps"
	case NalUnitTypeDps:
		return "dps"
	case NalUnitTypeReserved17:
		return "reserved 17"
	case NalUnitTypeReserved18:
		return "reserved 18"
	case NalUnitTypeAuxiliarySlice:
		return "auxiliary slice"
	case NalUnitTypeExtenSlice:
		return "exten slice"
	case NalUnitTypeDepthExtenSlice:
		return "depth exten slice"
	case NalUnitTypeReserved22:
		return "reserved 22"
	case NalUnitTypeReserved23:
		return "reserved 23"
	case NalUnitTypeUnspecified24:
		return "unspecified 24"
	case NalUnitTypeUnspecified25:
		return "unspecified 25"
	case NalUnitTypeUnspecified26:
		return "unspecified 26"
	case NalUnitTypeUnspecified27:
		return "unspecified 27"
	case NalUnitTypeUnspecified28:
		return "unspecified 28"
	case NalUnitTypeUnspecified29:
		return "unspecified 29"
	case NalUnitTypeUnspecified30: // 30
		return "unspecified 30"
	case NalUnitTypeUnspecified31: // 31
		return "unspecified 31"
	default:
		return fmt.Sprintf("unknown nalu type %d", naluType)
	}
}

func ParseNalUnitType(nalu []byte) NalUnitType {
	var naluType NalUnitType
	_ = (&naluType).Unmarshal(nalu[NalUnitTypeByteIndex:])
	return naluType
}
