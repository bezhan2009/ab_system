package uuid

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func NewUUID4() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}

	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:16],
	), nil
}

func ValidateUUID(id string) error {
	_, err := uuid.Parse(id)
	return err
}
