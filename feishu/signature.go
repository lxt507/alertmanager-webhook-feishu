package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

// VerifySign 验证飞书机器人签名
func VerifySign(secret string, timestamp int64, sign string) error {
	if secret == "" {
		return nil
	}

	now := time.Now().Unix()
	if absInt64(now-timestamp) > 300 {
		return fmt.Errorf("timestamp expired: %d", timestamp)
	}

	// message = timestamp + "\n"
	// key = secret
	expectedSign, err := GenSign(secret, timestamp)
	if err != nil {
		return fmt.Errorf("generate sign failed: %w", err)
	}

	// 对比签名
	if sign != expectedSign {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}

func GenSign(secret string, timestamp int64) (string, error) {
	// message = timestamp + "\n"
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n"

	// HMAC-SHA256，使用 secret 作为 key
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}

	// base64 encode
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
