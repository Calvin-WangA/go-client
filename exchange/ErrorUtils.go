package exchange

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const defaultErrMsg = "Error occur:"

type errorMethod interface {

	IsSucc() bool

	IsFail() bool
	// 输出错误信息
	ErrorPrintln(err error)

	setMessage(msg string)

	GetMessage() string

	GetErrCode() string
}

type ExchangeError struct {
	// 错误码
	errCode int
	// 错误信息
	errMsg string
	// 错误信息参数
	errParameters []string
}

func (error ExchangeError) IsSucc() bool {
	if error.errCode == 0 {
		return true
	}

	return false
}

func (error ExchangeError) IsFail() bool {
	if error.errCode != 0 {
		return true
	}

	return false
}

func newExchangeError(errCode int) ExchangeError {

	errMsg := replaceErrMsgParameters(errCode, nil)
	return ExchangeError{
		errCode:       errCode,
		errMsg:        errMsg,
		errParameters: nil,
	}
}

func newExchangeErrorByParams(errCode int, errParameters []string) ExchangeError {

	errMsg := replaceErrMsgParameters(errCode, errParameters)
	return ExchangeError{
		errCode:       errCode,
		errMsg:        errMsg,
		errParameters: nil,
	}
}

func (error ExchangeError) ErrorPrintln(err error) {
	// 打印错误堆栈西悉尼
	log.Println(GetErrorStackf(err, error.errMsg))
}

func (error ExchangeError) setMessage(msg string) {
	error.errMsg = msg
}

func (error ExchangeError) GetMessage() string{
	return error.errMsg
}

func (error ExchangeError) GetErrCode() int {
	return error.errCode
}

func replaceErrMsgParameters(errCode int, errParameters []string) string {
	errMsg := errmsgs[strconv.Itoa(errCode)]
	paramSize := strings.Count(errMsg, "{}")
	if len(errParameters) > 0 {
		for i := 0; i < paramSize; i++ {
			errMsg = strings.Replace(errMsg, "{}", errParameters[i], i + 1)
		}
	}

	return errMsg
}

/**
 * 用特定信息新创建一个 error
 */
func NewErrorf(format string, a ...interface{}) error {
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}

	errMsg := fmt.Sprintf("error occur, cause: %s \n\tat %s:%d", fmt.Sprintf(format, a...), f.Name(), line)
	return errors.New(errMsg)
}

/**
 * 用特定信息包装一个 error，使其包含代码堆栈信息
 * 如果 err 为空则返回空
 */
func WrapError(err error, wrapMsg string) error {
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}

	var wrapErr error = nil
	if err != nil {
		if wrapMsg == "" {
			wrapMsg = defaultErrMsg
		}
		errMsg := fmt.Sprintf("%s \n\tat %s:%d\nCause by: %s", wrapMsg, f.Name(), line, err.Error())
		wrapErr = errors.New(errMsg)
	}
	return wrapErr
}

/**
 * 用特定信息包装一个 error，使其包含代码堆栈信息
 * 如果 err 为空则返回空
 */
func WrapErrorf(err error, wrapMsgFmt string, a ...interface{}) error {
	wrapMsg := ""
	if wrapMsgFmt == "" {
		wrapMsg = defaultErrMsg
	} else {
		wrapMsg = fmt.Sprintf(wrapMsgFmt, a...)
	}

	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}

	var wrapErr error = nil
	if err != nil {
		if wrapMsg == "" {
			wrapMsg = defaultErrMsg
		}
		errMsg := fmt.Sprintf("%s \n\tat %s:%d\nCause by: %s", wrapMsg, f.Name(), line, err.Error())
		wrapErr = errors.New(errMsg)
	}
	return wrapErr
}

/**
 * 获得错误描述的同时携带上代码堆栈信息
 */
func GetErrorStack(err error, preStr string) string {
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return "GetErrorStack 方法获取堆栈失败，返回错误原信息：" + err.Error()
	}

	var errMsg string
	if err != nil{
		if preStr == "" {
			preStr = defaultErrMsg
		}
		errMsg = fmt.Sprintf("%s \n\tat %s:%d\nCause by: %s", preStr, f.Name(), line, err.Error())
	}
	return errMsg
}

/**
 * 获得错误描述的同时携带上代码堆栈信息
 */
func GetErrorStackf(err error, preStrFmt string, a ...interface{}) string {
	preStr := ""
	if preStrFmt == "" {
		preStr = defaultErrMsg
	} else {
		preStr = fmt.Sprintf(preStrFmt, a...)
	}
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return "GetErrorStack 方法获取堆栈失败，返回错误原信息：" + err.Error()
	}

	var errMsg string
	if err != nil{
		if preStr == "" {
			preStr = defaultErrMsg
		}
		errMsg = fmt.Sprintf("%s \n\tat %s:%d\nCause by: %s", preStr, f.Name(), line, err.Error())
	}
	return errMsg
}
