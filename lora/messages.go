package lora

import "github.com/brocaar/lorawan"

// GatewayMessage contains payloads received from lora gateway bridge
type GatewayMessage struct {
	RxInfo     RxInfo             `json:"rxInfo"`
	PhyPayload lorawan.PHYPayload `json:"phyPayload"`
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
