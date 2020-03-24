package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/I1820/gm/request"
	"github.com/labstack/echo/v4"
)

func (suite *Suite) TestDecrypt() {

	data, err := json.Marshal(request.Decrypt{
		AppSKey: "2B7E151628AED2A6ABF7158809CF4F3C",
		NetSKey: "2B7E151628AED2A6ABF7158809CF4F3C",
		PhyPayload: []byte{0x40, 0x30, 0x00, 0x00, 0x00, 0x00, 0xCC, 0x18, 0x01, 0x19,
			0xC8, 0x00, 0x1A, 0x8A, 0x2C, 0xAF, 0x60, 0x59, 0x8F, 0x17, 0x87, 0xCD, 0xDE, 0x2C, 0x6B, 0x43},
	})
	suite.NoError(err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/decrypt", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code, w.Body.String())
}
