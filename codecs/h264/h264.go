package h264

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
