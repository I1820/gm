package lora

import "fmt"

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
	// MTypeJoinAccept join accept MAC message Type
	MTypeJoinAccept
	// MTypeUnconfDataUp unconfirmed MAC message Type
	MTypeUnconfDataUp
	// MTypeUnconfDataDown unconfirmed data down MAC message Type
	MTypeUnconfDataDown
	// MTypeConfDataUp confirmed data up MAC message Type
	MTypeConfDataUp
	// MTypeConfDataDown confirmed data down MAC message Type
	MTypeConfDataDown
	// MTypeRejoinRequest rejoin request MAC message Type
	MTypeRejoinRequest
	// MTypeProprietary proprietary MAC message Type
	MTypeProprietary
)

// PhyPayload Physical Payload
type PhyPayload struct {
	MHDR
	FHDR
}

// UnmarshalJSON unmarshals physical payload binary data
func (p *PhyPayload) UnmarshalJSON(b []byte) error {
	p.MHDR.MType = (uint8(b[0]) & 0x07)
	p.MHDR.Major = (uint8(b[0]) & 0x38) >> 3
	p.MHDR.Major = (uint8(b[0]) & 0xc0) >> 6

	p.FHDR.DevAddr = fmt.Sprintf("%2x%2x%2x%2x", b[4], b[3], b[2], b[1])
	return nil
}

// MHDR MAC header specifices the message type (MType) and according to wich major version (Major)
// of the frame format of the LoRaWAN layer specification the frame has been encoded.
type MHDR struct {
	Major uint8
	RFU   uint8
	MType uint8
}

// FHDR ...
type FHDR struct {
	DevAddr string
}
