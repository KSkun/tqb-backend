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

type ReqUserAddUser struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RspUserAddUser struct {
	ID string `json:"_id"`
}

type ReqUserSendVerifyMail struct {
	Email string `query:"email" validate:"required,email"`
}

type ReqUserVerifyEmail struct {
	VerifyID string `query:"verify_id" validate:"required,uuid"`
}

type ReqUserChangePassword struct {
	VerifyID string `query:"verify_id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type RspUserGetInfo struct {
	ID               string   `json:"_id"`
	Username         string   `json:"username"`
	Email            string   `json:"email"`
	IsEmailVerified  bool     `json:"is_email_verified"`
	LastQuestion     string   `json:"last_question"`
	LastScene        string   `json:"last_scene"`
	StartTime        int64    `json:"start_time"`
	UnlockedScene    []string `json:"unlocked_scene"`
	FinishedQuestion []string `json:"finished_question"`
}

type ObjRspSubmissionQuestion struct {
	ID    string `json:"_id"`
	Title string `json:"title"`
}

type ObjRspSubmission struct {
	ID         string                   `json:"_id"`
	Time       int64                    `json:"time"`
	Question   ObjRspSubmissionQuestion `json:"question"`
	File       []string                 `json:"file"`
	Option     [][]int                  `json:"option"`
	Point      float64                  `json:"point"`
	AnswerTime int                      `json:"answer_time"`
}

type RspUserGetSubmission struct {
	Submission []ObjRspSubmission `json:"submission"`
}
