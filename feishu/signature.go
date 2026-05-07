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
	// 1. 将 timestamp + "\n" + 密钥 当做签名字符串
	// 2. 使用 HmacSHA256 算法计算空字符串的签名结果
	// 3. 再进行 Base64 编码
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	// 使用 stringToSign 作为 HMAC 的 key
	h := hmac.New(sha256.New, []byte(stringToSign))

	// 对空字符串进行 HMAC 计算
	_, err := h.Write([]byte(""))
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
