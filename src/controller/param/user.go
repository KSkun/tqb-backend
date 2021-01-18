package param

type ReqUserGetPublicKey struct {
	Email string `query:"email" validate:"required,email"`
}

type RspUserGetPublicKey struct {
	PublicKey string `json:"public_key"`
}
