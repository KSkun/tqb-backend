package param

const (
	StatusUnlock    = 0 // 未解锁
	StatusAnswering = 1 // 正在作答
	StatusFinish    = 2 // 已提交
)

type ObjRspQuestion struct {
	ID        string   `json:"_id"`
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	NextScene []string `json:"next_scene"`
	Status    int      `json:"status"`
}

type RspQuestionGetList struct {
	Question []ObjRspQuestion `json:"question"`
}

type ObjRspSubQuestion struct {
	Type      int      `json:"type"`
	Desc      string   `json:"desc"`
	Option    []string `json:"option"`
	FullPoint float64  `json:"full_point"`
	PartPoint float64  `json:"part_point"`
}

type ObjRspNextScene struct {
	Scene  string `json:"scene"`
	Option string `json:"option"`
}

type RspQuestionGetInfo struct {
	Title       string              `json:"title"`
	Desc        string              `json:"desc"`
	Statement   string              `json:"statement"`
	SubQuestion []ObjRspSubQuestion `json:"sub_question"`
	Author      string              `json:"author"`
	Audio       string              `json:"audio"`
	TimeLimit   int                 `json:"time_limit"`
	NextScene   []ObjRspNextScene   `json:"next_scene"`
	Status      int                 `json:"status"`
}

type ReqQuestionAddSubmission struct {
	Option [][]int  `json:"option"`
	File   []string `json:"file"`
}

type RspQuestionAddSubmission struct {
	ID string `json:"_id"`
}
