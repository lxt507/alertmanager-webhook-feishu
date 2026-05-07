package feishu

import (
	"testing"
	"time"
)

func TestVerifySign(t *testing.T) {
	secret := "test-secret-key-12345"
	now := time.Now().Unix()

	sign, err := GenSign(secret, now)
	if err != nil {
		t.Fatalf("GenSign failed: %v", err)
	}

	err = VerifySign(secret, now, sign)
	if err != nil {
		t.Errorf("VerifySign failed with valid signature: %v", err)
	}

	err = VerifySign(secret, now-600, sign)
	if err == nil {
		t.Error("VerifySign should fail with expired timestamp")
	}

	err = VerifySign(secret, now, "wrong-signature")
	if err == nil {
		t.Error("VerifySign should fail with wrong signature")
	}

	err = VerifySign("", now, "any-sign")
	if err != nil {
		t.Errorf("VerifySign should pass with empty secret: %v", err)
	}
}

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
