package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// EncryptAES AES-256-CBC加密
func EncryptAES(plaintext string) (string, error) {
	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %w", err)
	}

	plaintextBytes := []byte(plaintext)
	// PKCS7填充
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	for i := 0; i < padding; i++ {
		plaintextBytes = append(plaintextBytes, byte(padding))
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("生成IV失败: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintextBytes)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES AES-256-CBC解密
func DecryptAES(encrypted string) (string, error) {
	key := getEncryptionKey()
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建cipher失败: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("密文太短")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("密文长度不是块大小的倍数")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除PKCS7填充
	if len(ciphertext) > 0 {
		padding := int(ciphertext[len(ciphertext)-1])
		if padding > 0 && padding <= aes.BlockSize {
			ciphertext = ciphertext[:len(ciphertext)-padding]
		}
	}

	return string(ciphertext), nil
}

// MaskSecret 脱敏显示密钥（显示前4后4，中间用****代替）
func MaskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// getEncryptionKey 获取加密密钥（从环境变量读取，32字节用于AES-256）
func getEncryptionKey() []byte {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = "anchor-finance-default-key-32bytes!" // 默认密钥，生产环境应修改
	}
	// 确保密钥长度为32字节（AES-256）
	if len(key) < 32 {
		key = key + "00000000000000000000000000000000"[:32-len(key)]
	}
	return []byte(key[:32])
}
