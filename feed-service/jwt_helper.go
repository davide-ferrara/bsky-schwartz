package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "forzaReggina1914!forzaReggina1914!"
const aesKey = "12345678901234567890123456789012"

func GenerateToken(handle string, appPassword string) (string, error) {
	encryptedPassword, err := encryptPassword(appPassword)
	if err != nil {
		return "", fmt.Errorf("could not encrypt password: %w", err)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"handle":      handle,
		"appPassword": encryptedPassword,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign JWT: %w", err)
	}

	return tokenStr, nil
}

func ValidateToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("could not parse token: %w", err)
	}

	if !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid claims")
	}

	handle, ok := claims["handle"].(string)
	if !ok {
		return "", "", fmt.Errorf("handle not found in token")
	}

	encryptedPassword, ok := claims["appPassword"].(string)
	if !ok {
		return "", "", fmt.Errorf("appPassword not found in token")
	}

	appPassword, err := decryptPassword(encryptedPassword)
	if err != nil {
		return "", "", fmt.Errorf("could not decrypt password: %w", err)
	}

	return handle, appPassword, nil
}

func encryptPassword(password string) (string, error) {
	key := []byte(aesKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return hex.EncodeToString(ciphertext), nil
}

func decryptPassword(encrypted string) (string, error) {
	key := []byte(aesKey)
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
