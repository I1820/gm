/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 06-02-2018
 * |
 * | File Name:     runner/runner_test.go
 * +===============================================
 */

package gateway

import "testing"

func TestBasic(t *testing.T) {
	r, err := New("Eli", "127.0.0.1:1883")

	if err != nil {
		t.Fatalf("Runner creation error: %s", err)
	}

	t.Log(r.ID)

	err = r.Remove()

	if err != nil {
		t.Fatalf("Runner remove error: %s", err)
	}
}
