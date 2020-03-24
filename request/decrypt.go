package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Decrypt is the given payload with given keys
type Decrypt struct {
	PhyPayload []byte `json:"phy_payload"`
	AppSKey    string `json:"appskey"`
	NetSKey    string `json:"netskey"`
}

func (d Decrypt) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.PhyPayload, validation.Required),
		validation.Field(&d.AppSKey, validation.Required),
		validation.Field(&d.NetSKey, validation.Required),
	)
}
