package config

type ClaimToken struct {
	Secret string `env:"CLAIM_TOKEN_SECRET,required"`
}
