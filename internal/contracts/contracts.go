package contracts

import (
	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
)

const PROCESS = "lJLnoDsq8z0NJrTbQqFQ1arJayfuqWPqwRaW_3aNCgk"

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
