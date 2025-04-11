package utils

import (
	"crypto/rand"
	"fmt"
)

type UUID [16]byte

func NewUUID() (UUID, error) {
	var uuid UUID
	_, err := rand.Read(uuid[:])
	if err != nil {
		return UUID{}, fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Устанавливаем версию (4) в 13-й байт: 0100xxxx
	uuid[6] &= 0x0F // Очищаем первые 4 бита
	uuid[6] |= 0x40 // Устанавливаем версию 4

	// Устанавливаем вариант (10) в 17-й байт: 10xxxxxx
	uuid[8] &= 0x3F // Очищаем первые 2 бита
	uuid[8] |= 0x80 // Устанавливаем вариант 10

	return uuid, nil
}

// uuid в стр формат
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// парсит uuid из строки в байты
func ParseUUID(s string) (UUID, error) {
	var u UUID
	n, err := fmt.Sscanf(s, "%08x-%04x-%04x-%04x-%012x",
		&u[0], &u[4], &u[6], &u[8], &u[10])
	if err != nil || n != 5 {
		return UUID{}, fmt.Errorf("invalid UUID format: %s", s)
	}
	return u, nil
}

// чек если все байты равны нулю
func (u UUID) IsZero() bool {
	for _, b := range u[:] {
		if b != 0 {
			return false
		}
	}
	return true
}
