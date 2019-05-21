package hevc

import (
	"encoding/binary"
	"fmt"
)

//*    0                   1
//*    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
//*   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//*   |F|   Type    |  LayerId  | TID |
//*   +-------------+-----------------+
const (
	NalLayerIdMask      = 0x01f8
	NalLayerIdOffset    = 3
	NalLayerIdByteIndex = 0
)

type NalLayerId uint8

func (layerId NalLayerId) String() string {
	return fmt.Sprintf("%d (%x)", layerId, uint8(layerId))
}

func (layerId NalLayerId) Uint16() uint16 {
	b, _ := layerId.Marshal()
	return binary.BigEndian.Uint16(b)
}

func (layerId NalLayerId) Marshal() ([]byte, error) {
	var id uint16
	id = uint16(layerId) << NalLayerIdOffset
	var idBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(idBytes, id)
	return idBytes, nil
}

func (layerId *NalLayerId) Unmarshal(nalu []byte) error {

	*layerId = NalLayerId((binary.BigEndian.Uint16(nalu) & NalLayerIdMask) >> NalLayerIdOffset)
	return nil
}

func ParseNalLayerId(nalu []byte) NalLayerId {
	var layerId NalLayerId
	_ = (&layerId).Unmarshal(nalu[NalLayerIdByteIndex:])
	return layerId
}
