package lora

// GatewayMessage contains payloads received from lora gateway bridge
type GatewayMessage struct {
	RxInfo     RxInfo      `json:"rxInfo"`
	PhyPayload *PhyPayload `json:"phyPayload"`
}

// RxInfo contains gateway infomation that payloads
// received from it.
type RxInfo struct {
	Mac       string
	Timestamp int64
	RSSI      int     `json:"rssi"`
	LoRaSNR   float64 `json:"LoRaSNR"`
	Frequency int
	Size      int
	Channel   int
	CodeRate  string `json:"codeRate"`
}

const (
	// MTypeJoinRequest join request MAC message Type
	MTypeJoinRequest = iota
)

// PhyPayload ...
type PhyPayload struct {
	MHDR
	FHDR
}

// UnmarshalJSON ...
func (p *PhyPayload) UnmarshalJSON(b []byte) error {
	p.MHDR.MType = (uint8(b[0]) & 0x0e)
	return nil
}

// MHDR ...
type MHDR struct {
	Major uint8
	MType uint8
}

// FHDR ...
type FHDR struct {
	DevAddr uint32
}
