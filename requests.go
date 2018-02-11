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

// project request payload
type gatewayReq struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}
