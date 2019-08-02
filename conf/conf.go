package conf

var AUTH_USER string
var AUTH_PASS string

const (
	_ = iota
	E_CONTAINER_CREATE
	E_IMAGE_PULL
	S_CONTAINER_SUCCESS
)

type RespMsg struct {
	Msg  string `json:"msg,omitempty"`
	Code int    `json:"code,omitempty"`
}
