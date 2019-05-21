package hevc

const (
	// 7.4.3.1: vps_max_layers_minus1 is in [0, 62].
	MaxLayers = 63
	// 7.4.3.1: vps_max_sub_layers_minus1 is in [0, 6].
	MaxSubLayers = 7
	// 7.4.3.1: vps_num_layer_sets_minus1 is in [0, 1023].
	MaxLayerSets = 1024

	// 7.4.2.1: vps_video_parameter_set_id is u(4).
	MaxVpsCount = 16
	// 7.4.3.2.1: sps_seq_parameter_set_id is in [0, 15].
	MaxSpsCount = 16
	// 7.4.3.3.1: pps_pic_parameter_set_id is in [0, 63].
	MaxPpsCount = 64

	// A.4.2: MaxDpbSize is bounded above by 16.
	MaxDpbSize = 16
	// 7.4.3.1: vps_max_dec_pic_buffering_minus1[i] is in [0, MaxDpbSize - 1].
	MaxRefs = MaxDpbSize

	// 7.4.3.2.1: num_short_term_ref_pic_sets is in [0, 64].
	MaxShortTermRefPicSets = 64
	// 7.4.3.2.1: num_long_term_ref_pics_sps is in [0, 32].
	MaxLongTermRefPics = 32

	// A.3: all profiles require that CtbLog2SizeY is in [4, 6].
	MinLog2CtbSize = 4
	MaxLog2CtbSize = 6

	// E.3.2: cpb_cnt_minus1[i] is in [0, 31].
	MaxCpbCnt = 32

	// A.4.1: in table A.6 the highest level allows a MaxLumaPs of 35 651 584.
	MaxLumaPs = 35651584
	// A.4.1: pic_width_in_luma_samples and pic_height_in_luma_samples are
	// constrained to be not greater than sqrt(MaxLumaPs * 8).  Hence height/
	// width are bounded above by sqrt(8 * 35651584) = 16888.2 samples.
	MaxWidth  = 16888
	MaxHeight = 16888

	// A.4.1: table A.6 allows at most 22 tile rows for any level.
	MaxTileRows = 22
	// A.4.1: table A.6 allows at most 20 tile columns for any level.
	MaxTileColumns = 20

	// 7.4.7.1: in the worst case (tiles_enabled_flag and
	// entropy_coding_sync_enabled_flag are both set), entry points can be
	// placed at the beginning of every Ctb row in every tile, giving an
	// upper bound of (num_tile_columns_minus1 + 1) * PicHeightInCtbsY - 1.
	// Only a stream with very high resolution and perverse parameters could
	// get near that, though, so set a lower limit here with the maximum
	// possible value for 4K video (at most 135 16x16 Ctb rows).
	MaxEntryPointOffsets = MaxTileColumns * 135
)
