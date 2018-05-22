package lora

import "time"

// GatewayMessage contains payloads received from your nodes
type GatewayMessage struct {
	PhyPayload []byte
	RxInfo     RxInfo `json:"rxInfo"`
}

// RxInfo contains gateway infomation that payloads
// received from it.
type RxInfo struct {
	Mac       string
	Name      string
	Timestamp time.Time
	RSSI      int     `json:"rssi"`
	LoRaSNR   float64 `json:"LoRaSNR"`
	Frequency int64
	Size      int
	Channel   int
}
