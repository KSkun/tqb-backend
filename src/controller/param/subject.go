package param

type ObjRspSubject struct {
	Abbr       string `json:"abbr"`
	Name       string `json:"name"`
	StartScene string `json:"start_scene"`
}

type RspSubjectGetList struct {
	Subject []ObjRspSubject `json:"subject"`
}
