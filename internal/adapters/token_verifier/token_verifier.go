package tokenVerifier

import "slices"

type Verifier struct {
	Tokens []string
}

func (v *Verifier) VerifyToken(token string) bool {
	return slices.Contains(v.Tokens, token)
}
