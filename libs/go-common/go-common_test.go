package gocommon

import (
	"testing"
)

func TestGoCommon(t *testing.T) {
	result := GoCommon("works")
	if result != "GoCommon works" {
		t.Error("Expected GoCommon to append 'works'")
	}
}
