package base

import (
	"github.com/pquerna/ffjson/ffjson"
	"launchpad/constant/serror"
)

type Response struct {
	Data any    `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *Response) InitCode(code int) {
	r.Code = code
	r.Msg = serror.ErrMsg(code)
}

func (r *Response) ToStruct(obj any) error {
	b, _ := ffjson.Marshal(r.Data)
	return ffjson.Unmarshal(b, obj)
}
