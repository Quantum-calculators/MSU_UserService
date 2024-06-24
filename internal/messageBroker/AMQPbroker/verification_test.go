package messagebroker_test

import (
	"testing"

	messagebroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker/AMQPbroker"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	ch, err := messagebroker.TestBroker()
	assert.NoError(t, err)
	broker := messagebroker.New(ch)

	err = broker.Message().SendMessage([]byte("test message"), "/VerifyEmail")
	assert.NoError(t, err)
}
