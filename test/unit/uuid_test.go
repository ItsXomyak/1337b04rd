package unit

import (
	"testing"

	common "1337b04rd/internal/app/common/utils"
)

func TestNewUUID(t *testing.T) {
	id1, err := common.NewUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}
	id2, err := common.NewUUID()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}
	if id1 == id2 {
		t.Error("UUIDs should be unique")
	}
}