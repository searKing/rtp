package hevc

import "fmt"

//*    0                   1
//*    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
//*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//*   |F|   Type    |  LayerId  | TID |
//*   +-------------+-----------------+
const (
	NalTemporalIdMask      = 0x7
	NalTemporalIdOffset    = 0
	NalTemporalIdByteIndex = 1
)

type NalTemporalId uint8

func (tid NalTemporalId) Byte() byte {
	b, _ := tid.Marshal()
	return b[0]
}

func (tid NalTemporalId) Marshal() ([]byte, error) {
	return []byte{byte(tid) << NalTemporalIdOffset}, nil
}

func (tid *NalTemporalId) Unmarshal(buf []byte) error {
	*tid = NalTemporalId((buf[0] & NalTemporalIdMask) >> NalTemporalIdOffset)
	return nil
}

func ParseNalTemporalId(nalu []byte) NalTemporalId {
	var tid NalTemporalId
	_ = (&tid).Unmarshal(nalu[NalTemporalIdByteIndex:])
	return tid
}

func (tid NalTemporalId) String() string {
	return fmt.Sprintf("%d (%x)", tid, uint8(tid))
}
