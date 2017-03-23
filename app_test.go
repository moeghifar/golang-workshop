package main

import "testing"

func TestRedisPing(t *testing.T) {
	var resVal string
	resVal, _ = redisPing()
	if resVal != "PONG" {
		t.Error("You've got ", resVal, "-> Expected :: PONG")
	}
}
