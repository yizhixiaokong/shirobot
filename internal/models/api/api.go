package api

// Request API请求结构体
type Request struct {
	Action string `json:"action"` // 终结点名称, 例如 'send_group_msg'
	Params any    `json:"params"` // 参数
	Echo   string `json:"echo"`   // 如果指定了 echo 字段, 那么响应包也会同时包含一个 echo 字段, 它们会有相同的值
}

// NewRequest 新建API请求结构体
func NewRequest(action string, params any, echo string) Request {
	// TODO:  参数Option化
	return Request{
		Action: action,
		Params: params,
		Echo:   echo,
	}
}

// Response API请求返回体
type Response struct {
	Status  string          `json:"status"`  // 状态, 表示 API 是否调用成功, 如果成功, 则是 OK, 其他的在下面会说明
	Retcode ResponseRetcode `json:"retcode"` // 错误码
	Msg     string          `json:"msg"`     // 错误消息, 仅在 API 调用失败时有该字段
	Wording string          `json:"wording"` // 对错误的详细解释(中文), 仅在 API 调用失败时有该字段
	Data    any             `json:"data"`    // 相应返回数据
	Echo    string          `json:"echo"`    // 如果请求时指定了 echo, 那么响应也会包含 echo
}

// ResponseRetcode 错误码
type ResponseRetcode int

// ok	api 调用成功
// async	api 调用已经提交异步处理, 此时 retcode 为 1, 具体 api 调用是否成功无法得知
// failed	api 调用失败 (操作失败, 具体原因可以看响应的 msg 和 wording 字段)
const (
	OK ResponseRetcode = iota
	ASYNC
	FAILED
)

func (rs ResponseRetcode) String() string {
	switch rs {
	case OK:
		return "ok"
	case ASYNC:
		return "async"
	case FAILED:
		return "failed"
	default:
		return "unknow"
	}
}
