package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenSign 生成签名,用于签名校验
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
