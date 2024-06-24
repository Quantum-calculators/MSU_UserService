package testbroker_test

import (
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/messageBroker/testbroker"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	broker := testbroker.New()

	err := broker.Message().SendMessage([]byte("test message"), "/VerifyEmail")
	assert.NoError(t, err)
}
