package gapi

import (
	"fmt"

	"github.com/hkloudou/xlib/xcolor"
)

type appErrorSource string

const (
	appErrorSourceSign     appErrorSource = "API鉴权"
	appErrorSourceValidate appErrorSource = "数据校验"
)

type appError struct {
	code   int
	source appErrorSource
	msg    string
	block  string
	raw    error
}

func (e appError) Error() string {
	return fmt.Sprintf("%s错误：%s", e.source, e.msg)
}

func (e appError) String() string {
	str := fmt.Sprintf("[%10s] %s %s", xcolor.Red(string(e.source)), xcolor.Blue(e.block), e.msg)
	if e.raw != nil {
		str += fmt.Sprintf(" RAW: %s", xcolor.Yellow(e.raw.Error()))
	}
	return str
}

func NewAppErrorValidate(block, msg string, raw ...error) appError {
	err := appError{
		code:   1004,
		source: appErrorSourceValidate,
		block:  block,
		msg:    msg,
	}
	if len(raw) > 0 {
		err.raw = raw[0]
	}
	return err
}

func NewAppErrorSign(block, msg string, raw ...error) appError {
	err := appError{
		code:   1003,
		source: appErrorSourceSign,
		block:  block,
		msg:    msg,
	}
	if len(raw) > 0 {
		err.raw = raw[0]
	}
	return err
}
