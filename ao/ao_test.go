package ao

import (
	"testing"

	"github.com/liteseed/bungo/internal/types"
	"gotest.tools/v3/assert"
)

func TestSendMessage(t *testing.T) {
	args := SendMessageArgs{Data: "", Tags: []types.Tag{{Name: "Action", Value: "Staker"}}}
	mId, err := SendMessage(args)
	assert.NilError(t, err)
	t.Log(mId)
}

func TestReadResult(t *testing.T) {
	args := ReadResultArgs{Message: "zq-lHp0drHxKi-GrqP4EIs5qrMXRT-4mjfiLikxCZ0"}
	r, err := ReadResult(args)
	assert.NilError(t, err)
	t.Log(r.GasUsed)
}
