package unit

import (
	"1337b04rd/internal/app/common/utils"
	"testing"
<<<<<<< HEAD

	common "1337b04rd/internal/app/common/utils"
=======
>>>>>>> a18a7b86809f5a800966d699c87607c95a839569
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
