package unit

import (
	"1337b04rd/internal/app/common/utils"
	"testing"
)

func TestNewUUID(t *testing.T) {
	id1, err := utils.NewUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}
	id2, err := utils.NewUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}
	if id1 == id2 {
		t.Error("UUIDs should be unique")
	}
}
