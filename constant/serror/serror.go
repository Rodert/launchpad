package serror

type Error struct {
	err  string
	Code int
}

func (m *Error) Error() string {
	return m.err
}

func New(code int) error {
	if code == Success {
		return nil
	}
	return &Error{
		err:  ErrMsg(code),
		Code: code,
	}
}

const (
	Success         = 0
	Unknown         = 99999
	CallServiceFail = 10001
	ParamsErr       = 10002
	ServiceNotExist = 20001
	DBConnectErr    = 30001
	RedisConnectErr = 40001
	CustomErrCode   = 50001
	NotFoundErrCode = 60001
	SystemErrCode   = 70001
)

const (
	CallServiceFailMsg = "CallServiceFail"
	UnknownMsg         = "Unknown"
	SuccessMsg         = "Success"
	ServiceNotExistMsg = "ServiceNotExist"
	DBConnectErrMsg    = "DBConnectErr"
	RedisConnectErrMsg = "RedisConnectErr"
	ParamsErrMsg       = "ParamsErr"
	CustomErrMsg       = "CustomErr"
	NotFoundErrMsg     = "NotFoundErr"
	SystemErrMsg       = "SystemErr"
)

var ErrMp = map[int]string{
	Success:         SuccessMsg,
	ParamsErr:       ParamsErrMsg,
	CallServiceFail: CallServiceFailMsg,
	ServiceNotExist: ServiceNotExistMsg,
	Unknown:         UnknownMsg,
	DBConnectErr:    DBConnectErrMsg,
	RedisConnectErr: RedisConnectErrMsg,
	CustomErrCode:   CustomErrMsg,
	NotFoundErrCode: NotFoundErrMsg,
	SystemErrCode:   SystemErrMsg,
}

func ErrMsg(code int) string {
	if _, ok := ErrMp[code]; ok {
		return ErrMp[code]
	}
	return UnknownMsg
}

// CustomErr 业务自定义 error 非特殊需求 code默认传0
func CustomErr(msg string, code int) error {
	if code == 0 {
		code = CustomErrCode
	}
	return &Error{
		err:  msg,
		Code: code,
	}
}
