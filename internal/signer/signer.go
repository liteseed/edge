package signer

import "github.com/everFinance/goar"

type Signer struct {
	goarSigner *goar.Signer
}

func New(path string) *Signer {
	s, err := goar.NewSignerFromPath(path)
	if err != nil {
		return nil
	}
	return &Signer{s}
}

func (s *Signer) sign(data []byte) ([]byte, error) {
	return s.goarSigner.SignMsg(data)
}

func (s *Signer) address() string {
	return s.goarSigner.Address
}
