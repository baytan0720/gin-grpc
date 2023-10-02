package gin_grpc

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	UUID   string
	UserId int64
	Force  bool
}

func TestStoreRequestIntoKeys(t *testing.T) {
	f := StoreRequestIntoKeys()
	expected := &TestRequest{
		UUID:   uuid.New().String(),
		UserId: 1,
		Force:  true,
	}
	c := &Context{
		Req: expected,
	}
	f(c)
	assert.Equal(t, expected.UUID, c.GetString("UUID"))
	assert.Equal(t, expected.UserId, c.GetInt64("UserId"))
	assert.Equal(t, expected.Force, c.GetBool("Force"))
}
