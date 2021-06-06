package kmp

import (
	"testing"
)

func TestKmp_Compare(t *testing.T) {
	k := New("abcabcd")

	res := k.Compare("aaaabcabcdef")
	if res == -1 {
		t.Fatalf("exceped true, acutal false")
	}
}
