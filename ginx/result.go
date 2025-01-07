package ginx

type Result struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
