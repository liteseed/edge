package contracts

import (
	"github.com/everFinance/goar"
	"github.com/liteseed/aogo"
)

type Context struct {
	ao      *aogo.AO
	process string
	signer  *goar.ItemSigner
}

func New(ao *aogo.AO, process string, signer *goar.ItemSigner) *Context {
	return &Context{
		ao:      ao,
		process: process,
		signer:  signer,
	}
}
