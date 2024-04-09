package responsor

type Response struct {
	Domain  string   `json:"domain"`
	BizCode int64    `json:"biz_code"`
	Message string   `json:"message"`
	Reasons []string `json:"reasons"`
	Data    any      `json:"data"`
}
