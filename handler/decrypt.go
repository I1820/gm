package handler

import (
	"encoding/hex"
	"net/http"

	"github.com/I1820/gm/request"
	"github.com/brocaar/lorawan"
	"github.com/labstack/echo/v4"
)

type LoRa struct {
}

func (l LoRa) Register(g *echo.Group) {
	g.POST("/decrypt", l.Decrypt)
}

func (LoRa) Decrypt(c echo.Context) error {
	var rq request.Decrypt
	if err := c.Bind(&rq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	appSKeySlice, err := hex.DecodeString(rq.AppSKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var appSKey lorawan.AES128Key

	copy(appSKey[:], appSKeySlice)

	netSKeySlice, err := hex.DecodeString(rq.NetSKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var netSKey lorawan.AES128Key

	copy(netSKey[:], netSKeySlice)

	var phy lorawan.PHYPayload
	if err := phy.UnmarshalBinary(rq.PhyPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	mac, ok := phy.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "*MACPayload expected")
	}

	success, err := phy.ValidateUplinkJoinMIC(netSKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !success {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid MIC")
	}

	if err := phy.DecryptFRMPayload(appSKey); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	data, ok := mac.FRMPayload[0].(*lorawan.DataPayload)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "*DataPayload expected")
	}

	return c.JSON(http.StatusOK, data.Bytes)
}
