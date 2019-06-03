/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 05-06-2018
 * |
 * | File Name:     main_test.go
 * +===============================================
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbout(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api/about", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	if string(body) != "18.20 is leaving us" {
		t.Fatalf("who leaving us?! %q", body)
	}
}

func TestDecrypt(t *testing.T) {
	h := handle()
	s := httptest.NewServer(h)
	defer s.Close()

	d, err := json.Marshal(decryptReq{
		AppSKey: "2B7E151628AED2A6ABF7158809CF4F3C",
		NetSKey: "2B7E151628AED2A6ABF7158809CF4F3C",
		PhyPayload: []byte{0x40, 0x30, 0x00, 0x00, 0x00, 0x00, 0xCC, 0x18, 0x01, 0x19,
			0xC8, 0x00, 0x1A, 0x8A, 0x2C, 0xAF, 0x60, 0x59, 0x8F, 0x17, 0x87, 0xCD, 0xDE, 0x2C, 0x6B, 0x43},
	})
	if err != nil {
		t.Fatalf("Project request marshaling: %s\n", d)
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/decrypt", s.URL), "application/json", bytes.NewBuffer(d))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log(body)
}
