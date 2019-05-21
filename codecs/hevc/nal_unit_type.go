package hevc

import "fmt"

//   0                   1
//   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |F|   Type    |  LayerId  | TID |
//	+-------------+-----------------+
const (
	NalUnitTypeMask      = 0x3f << NalUnitTypeOffset
	NalUnitTypeOffset    = 1
	NalUnitTypeByteIndex = 0
)

// ffmpeg/libavcodec/hevc.h
/*
 * Table 7-1 – NAL unit type codes and NAL unit type classes
 * T-REC-H.265-201802-I!!PDF-E.pdf
 */
type NalUnitType uint8

const (
	NalUnitTypeTrailN       NalUnitType = iota // 0
	NalUnitTypeTrailR                          // 1
	NalUnitTypeTsaN                            // 2
	NalUnitTypeTsaR                            // 3
	NalUnitTypeStsaN                           // 4
	NalUnitTypeStsaR                           // 5
	NalUnitTypeRadlN                           // 6
	NalUnitTypeRadlR                           // 7
	NalUnitTypeRaslN                           // 8
	NalUnitTypeRaslR                           // 9
	NalUnitTypeRsvVclN10                       // 10
	NalUnitTypeRsvVclR11                       // 11
	NalUnitTypeRsvVclN12                       // 12
	NalUnitTypeRsvVclR13                       // 12
	NalUnitTypeRsvVclN14                       // 14
	NalUnitTypeRsvVclR15                       // 15
	NalUnitTypeBlaWLp                          // 16
	NalUnitTypeBlaWRadl                        // 17
	NalUnitTypeBlaNLp                          // 18
	NalUnitTypeIdrWRadl                        // 19
	NalUnitTypeIdrNLp                          // 20
	NalUnitTypeCraNut                          // 21
	NalUnitTypeRsvIrapVcl22                    // 22
	NalUnitTypeRsvIrapVcl23                    // 23
	NalUnitTypeRsvVcl24                        // 24
	NalUnitTypeRsvVcl25                        // 25
	NalUnitTypeRsvVcl26                        // 26
	NalUnitTypeRsvVcl27                        // 27
	NalUnitTypeRsvVcl28                        // 28
	NalUnitTypeRsvVcl29                        // 29
	NalUnitTypeRsvVcl30                        // 30
	NalUnitTypeRsvVcl31                        // 31
	NalUnitTypeVpsNut                          // 32
	NalUnitTypeSpsNut                          // 33
	NalUnitTypePpsNut                          // 34
	NalUnitTypeAudNut                          // 35
	NalUnitTypeEosNut                          // 36
	NalUnitTypeEobNut                          // 37
	NalUnitTypeFdNut                           // 38
	NalUnitTypePrefixSeiNut                    // 39
	NalUnitTypeSuffixSeiNut                    // 40
	NalUnitTypeRsvNvcl41                       // 41
	NalUnitTypeRsvNvcl42                       // 42
	NalUnitTypeRsvNvcl43                       // 43
	NalUnitTypeRsvNvcl44                       // 44
	NalUnitTypeRsvNvcl45                       // 45
	NalUnitTypeRsvNvcl46                       // 46
	NalUnitTypeRsvNvcl47                       // 47
	NalUnitTypeUnspec48                        // 48
	NalUnitTypeUnspec49                        // 49
	NalUnitTypeUnspec50                        // 50
	NalUnitTypeUnspec51                        // 51
	NalUnitTypeUnspec52                        // 52
	NalUnitTypeUnspec53                        // 53
	NalUnitTypeUnspec54                        // 54
	NalUnitTypeUnspec55                        // 55
	NalUnitTypeUnspec56                        // 56
	NalUnitTypeUnspec57                        // 57
	NalUnitTypeUnspec58                        // 58
	NalUnitTypeUnspec59                        // 59
	NalUnitTypeUnspec60                        // 60
	NalUnitTypeUnspec61                        // 61
	NalUnitTypeUnspec62                        // 62
	NalUnitTypeUnspec63                        // 63
)

func (naluType NalUnitType) Vcl() bool {
	return naluType < 32
}
func (naluType NalUnitType) Byte() byte {
	b, _ := naluType.Marshal()
	return b[0]
}

func (naluType NalUnitType) Marshal() ([]byte, error) {
	return []byte{byte(naluType) << NalUnitTypeOffset}, nil
}

func (naluType *NalUnitType) Unmarshal(buf []byte) error {
	*naluType = NalUnitType(buf[0] & NalUnitTypeMask >> NalUnitTypeOffset)
	return nil
}

func ParseNalUnitType(nalu []byte) NalUnitType {
	var naluType NalUnitType
	_ = (&naluType).Unmarshal(nalu[NalUnitTypeByteIndex:])
	return naluType
}

func (naluType NalUnitType) String() string {
	switch naluType {
	case NalUnitTypeTrailN:
		return "TRAIL_N"
	case NalUnitTypeTrailR:
		return "TRAIL_R"
	case NalUnitTypeTsaN:
		return "TSA_N"
	case NalUnitTypeTsaR:
		return "TSA_R"
	case NalUnitTypeStsaN:
		return "STSA_N"
	case NalUnitTypeStsaR:
		return "STSA_R"
	case NalUnitTypeRadlN:
		return "RADL_N"
	case NalUnitTypeRadlR:
		return "RADL_R"
	case NalUnitTypeRaslN:
		return "RASL_N"
	case NalUnitTypeRaslR:
		return "RASL_R"
	case NalUnitTypeRsvVclN10:
		return "RSV_VCL_N10"
	case NalUnitTypeRsvVclR11:
		return "RSV_VCL_R11"
	case NalUnitTypeRsvVclN12:
		return "RSV_VCL_N12"
	case NalUnitTypeRsvVclR13:
		return "RSV_VCL_R13"
	case NalUnitTypeRsvVclN14:
		return "RSV_VCL_N14"
	case NalUnitTypeRsvVclR15:
		return "RSV_VCL_R15"
	case NalUnitTypeBlaWLp:
		return "BLA_W_LP"
	case NalUnitTypeBlaWRadl:
		return "BLA_W_RADL"
	case NalUnitTypeBlaNLp:
		return "BLA_N_LP"
	case NalUnitTypeIdrWRadl:
		return "IDR_W_RADL"
	case NalUnitTypeIdrNLp:
		return "IDR_N_LP"
	case NalUnitTypeCraNut:
		return "CRA_NUT"
	case NalUnitTypeRsvIrapVcl22:
		return "RSV_IRAP_VCL22"
	case NalUnitTypeRsvIrapVcl23:
		return "RSV_IRAP_VCL23"
	case NalUnitTypeRsvVcl24:
		return "RSV_VCL24"
	case NalUnitTypeRsvVcl25:
		return "RSV_VCL25"
	case NalUnitTypeRsvVcl26:
		return "RSV_VCL26"
	case NalUnitTypeRsvVcl27:
		return "RSV_VCL27"
	case NalUnitTypeRsvVcl28:
		return "RSV_VCL28"
	case NalUnitTypeRsvVcl29:
		return "RSV_VCL29"
	case NalUnitTypeRsvVcl30:
		return "RSV_VCL30"
	case NalUnitTypeRsvVcl31:
		return "RSV_VCL31"
	case NalUnitTypeVpsNut:
		return "VPS_NUT"
	case NalUnitTypeSpsNut:
		return "SPS_NUT"
	case NalUnitTypePpsNut:
		return "PPS_NUT"
	case NalUnitTypeAudNut:
		return "AUD_NUT"
	case NalUnitTypeEosNut:
		return "EOS_NUT"
	case NalUnitTypeEobNut:
		return "EOB_NUT"
	case NalUnitTypeFdNut:
		return "FD_NUT"
	case NalUnitTypePrefixSeiNut:
		return "PREFIX_SEI_NUT"
	case NalUnitTypeSuffixSeiNut:
		return "SUFFIX_SEI_NUT"
	case NalUnitTypeRsvNvcl41:
		return "RSV_NVCL41"
	case NalUnitTypeRsvNvcl42:
		return "RSV_NVCL42"
	case NalUnitTypeRsvNvcl43:
		return "RSV_NVCL43"
	case NalUnitTypeRsvNvcl44:
		return "RSV_NVCL44"
	case NalUnitTypeRsvNvcl45:
		return "RSV_NVCL45"
	case NalUnitTypeRsvNvcl46:
		return "RSV_NVCL46"
	case NalUnitTypeRsvNvcl47:
		return "RSV_NVCL47"
	case NalUnitTypeUnspec48:
		return "UNSPEC48"
	case NalUnitTypeUnspec49:
		return "UNSPEC49"
	case NalUnitTypeUnspec50:
		return "UNSPEC50"
	case NalUnitTypeUnspec51:
		return "UNSPEC51"
	case NalUnitTypeUnspec52:
		return "UNSPEC52"
	case NalUnitTypeUnspec53:
		return "UNSPEC53"
	case NalUnitTypeUnspec54:
		return "UNSPEC54"
	case NalUnitTypeUnspec55:
		return "UNSPEC55"
	case NalUnitTypeUnspec56:
		return "UNSPEC56"
	case NalUnitTypeUnspec57:
		return "UNSPEC57"
	case NalUnitTypeUnspec58:
		return "UNSPEC58"
	case NalUnitTypeUnspec59:
		return "UNSPEC59"
	case NalUnitTypeUnspec60:
		return "UNSPEC60"
	case NalUnitTypeUnspec61:
		return "UNSPEC61"
	case NalUnitTypeUnspec62:
		return "UNSPEC62"
	case NalUnitTypeUnspec63:
		return "UNSPEC63"
	default:
		return fmt.Sprintf("unknown nalu type %d", naluType)
	}
}
