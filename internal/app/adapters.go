package app

import tokenVerifier "newsapi/internal/adapters/token_verifier"

type adapters struct {
	TokenVerifier *tokenVerifier.Verifier
}

func initAdapters(cfg Config) *adapters {
	return &adapters{
		TokenVerifier: &tokenVerifier.Verifier{
			Tokens: cfg.AuthTokens,
		},
	}
}
