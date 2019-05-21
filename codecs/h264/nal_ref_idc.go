package h264

const (
	NalRefIdcMask      = 0x60
	NalRefIdcOffset    = 5
	NalRefIdcByteIndex = 0
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

func (refId *NalRefIdc) Unmarshal(buf []byte) error {
	*refId = NalRefIdc((buf[0] & NalRefIdcMask) >> NalRefIdcOffset)
	return nil
}

func ParseNalRefIdc(nalu []byte) NalRefIdc {
	var refId NalRefIdc
	_ = (&refId).Unmarshal(nalu[NalRefIdcByteIndex:])
	return refId
}
