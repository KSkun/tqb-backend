package param

type ReqUserGetPublicKey struct {
	Email string `query:"email" validate:"required,email"`
}

type RspUserGetPublicKey struct {
	PublicKey string `json:"public_key"`
}

type ReqUserGetToken struct {
	Email    string `query:"email" validate:"required,email"`
	Password string `query:"password" validate:"required"`
}

type RspUserGetToken struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
