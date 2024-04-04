package contracts

import (
	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
)

const PROCESS = "jHohSHfn4t5_LXLDra4ApC9KhmexMjbVmCco-9jS0cQ"

type Context struct {
	ao     *aogo.AO
	signer *goar.ItemSigner
}

func New(ao *aogo.AO, signer *goar.ItemSigner) *Context {
	return &Context{
		ao:     ao,
		signer: signer,
	}
}
