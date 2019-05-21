package h264

const (
	ForbiddenZeroBitMask      = 1 << ForbiddenZeroBitOffset
	ForbiddenZeroBitOffset    = 7
	ForbiddenZeroBitByteIndex = 0
)

type ForbiddenZeroBit bool

func (f ForbiddenZeroBit) Byte() byte {
	b, _ := f.Marshal()
	return b[0]
}
func (f ForbiddenZeroBit) Marshal() ([]byte, error) {
	if f {
		return []byte{0}, nil
	}
	return []byte{ForbiddenZeroBitMask}, nil
}

func (f *ForbiddenZeroBit) Unmarshal(buf []byte) error {
	*f = ForbiddenZeroBit((buf[0] & ForbiddenZeroBitMask) != 0)
	return nil
}

func ParseForbiddenZeroBit(nalu []byte) ForbiddenZeroBit {
	var f ForbiddenZeroBit
	_ = (&f).Unmarshal(nalu[ForbiddenZeroBitByteIndex:])
	return f
}
