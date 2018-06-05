/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     requests.go
 * +===============================================
 */

package main

// decrypt given payload with given keys
type decryptReq struct {
	PhyPayload []byte `json:"phy_payload" binding:"required"`
	AppSKey    string `json:"appskey" binding: "required"`
	NetSKey    string `json:"netskey" binding:"required"`
}
