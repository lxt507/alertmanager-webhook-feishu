package feishu

import (
	"testing"
)

func TestGenSign(t *testing.T) {
	secret := "test-secret"
	timestamp := int64(1234567890)

	sign, err := GenSign(secret, timestamp)
	if err != nil {
		t.Fatalf("GenSign failed: %v", err)
	}

	if sign == "" {
		t.Error("GenSign should return non-empty signature")
	}

	sign2, _ := GenSign(secret, timestamp)
	if sign != sign2 {
		t.Error("GenSign should be deterministic")
	}
}
