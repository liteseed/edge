package signer

type Signer struct {
	address string
}

func (s *Signer) getAddress() string {
	return s.address
}
