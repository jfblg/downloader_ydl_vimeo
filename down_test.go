package main

import (
	"testing"
)

func TestTransFile(t *testing.T) {
	name := "1.3 Hips And Hamstrings D1"
	exp := "yoga15_1-3_Hips_And_Hamstrings_D1.mp4"
	act, _ := transFile(name)

	if act != exp {
		t.Errorf("Expected: %v, Got: %v\n", exp, act)
	}
}

func TestTransURL(t *testing.T) {
	url := "https://skyfire.vimeocdn.com/1165038331/master.json?base64_init=1"
	exp := "https://skyfire.vimeocdn.com/1165038331/master.mpd"

	act, _ := transURL(url)

	if act != exp {
		t.Errorf("\nExpected: %v\nGot: %v\n", exp, act)
	}
}
