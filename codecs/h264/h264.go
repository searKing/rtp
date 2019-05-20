package h264

const (
	NalUnitTypeMask = 0x1f

	NalRefIdcMask   = 0x60
	NalRefIdcOffset = 5
)

// ffmpeg/libavcodec/h264.h
/*
 * Table 7-1 – NAL unit type codes, syntax element categories, and NAL unit type classes in
 * T-REC-H.264-201704
 */
type NalUnitType uint8

func (naluType NalUnitType) Byte() byte {
	b, _ := naluType.Marshal()
	return b[0]
}
func (naluType NalUnitType) Marshal() ([]byte, error) {
	return []byte{byte(naluType)}, nil
}

func (naluType *NalUnitType) Unmarshal(nalu []byte) error {
	*naluType = NalUnitType(nalu[0] & NalUnitTypeMask)
	return nil
}

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

type NalRefIdc uint8

const (
	NalRefIdcNonIDRCodedSlice        NalRefIdc = 0x2
	NalRefIdcCodeSliceSataPartitionA NalRefIdc = 0x2
	NalRefIdcCodeSliceSataPartitionB NalRefIdc = 0x1
	NalRefIdcCodeSliceSataPartitionC NalRefIdc = 0x1
)

func (refId NalRefIdc) Byte() byte {
	b, _ := refId.Marshal()
	return b[0]
}
func (refId NalRefIdc) Marshal() ([]byte, error) {
	return []byte{byte(refId << NalRefIdcOffset)}, nil
}

func (refId *NalRefIdc) Unmarshal(nalu []byte) error {
	*refId = NalRefIdc(nalu[0]&NalRefIdcMask) >> NalRefIdcOffset
	return nil
}

func ParseNalUnitType(nalu []byte) NalUnitType {
	return NalUnitType(nalu[0] & NalUnitTypeMask)
}

func ParseNalRefIdc(nalu []byte) NalRefIdc {
	return NalRefIdc(nalu[0]&NalRefIdcMask) >> NalRefIdcOffset
}

var (
	StartSequence = []byte{0, 0, 0, 1}
)

const (
	// 7.4.2.1.1: seq_parameter_set_id is in [0, 31].
	MaxSpsCount = 32
	// 7.4.2.2: pic_parameter_set_id is in [0, 255].
	MaxPpsCount = 256

	// A.3: MaxDpbFrames is bounded above by 16.
	MaxDpbFrames = 16
	// 7.4.2.1.1: max_num_ref_frames is in [0, MaxDpbFrames], and
	// each reference frame can have two fields.
	MaxRefs = 2 * MaxDpbFrames

	// 7.4.3.1: modification_of_pic_nums_idc is not equal to 3 at most
	// num_ref_idx_lN_active_minus1 + 1 times (that is, once for each
	// possible reference), then equal to 3 once.
	MaxRplmCount = MaxRefs + 1

	// 7.4.3.3: in the worst case, we begin with a full short-term
	// reference picture list.  Each picture in turn is moved to the
	// long-term list (type 3) and then discarded from there (type 2).
	// Then, we set the length of the long-term list (type 4), mark
	// the current picture as long-term (type 6) and terminate the
	// process (type 0).
	MaxMmcoCount = MaxRefs*2 + 3

	// A.2.1, A.2.3: profiles supporting FMO constrain
	// num_slice_groups_minus1 to be in [0, 7].
	MaxSliceGroups = 8

	// E.2.2: cpb_cnt_minus1 is in [0, 31].
	MaxCpbCnt = 32

	// A.3: in table A-1 the highest level allows a MaxFS of 139264.
	MaxMbPicSize = 139264
	// A.3.1, A.3.2: PicWidthInMbs and PicHeightInMbs are constrained
	// to be not greater than sqrt(MaxFS * 8).  Hence height/width are
	// bounded above by sqrt(139264 * 8) = 1055.5 macroblocks.
	MaxMbWidth  = 1055
	MaxMbHeight = 1055
	MaxWidth    = MaxMbWidth * 16
	MaxHeight   = MaxMbHeight * 16
)
