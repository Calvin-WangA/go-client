package exchange

import (
	"bytes"
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
	// 返回参数
	filePath string
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

	var errMsg string
	if errCode != 0 {
		errMsg = replaceErrMsgParameters(errCode, nil)
	}
	return ExchangeError{
		errCode:       errCode,
		errMsg:        errMsg,
		errParameters: nil,
	}
}

func newExchangeErrorByParams(errCode int, errParameters []string) ExchangeError {

	var errMsg string
	if errCode != 0 {
		errMsg = replaceErrMsgParameters(errCode, errParameters)
	}
	return ExchangeError{
		errCode:       errCode,
		errMsg:        errMsg,
		errParameters: nil,
	}
}

func (error ExchangeError) ErrorPrintln(err error) {
	// 打印错误堆栈西悉尼
	if err == nil {
		log.Println("错误信息：" + error.errMsg)
	} else {
		log.Println(GetErrorStack(err, error.errMsg))
	}
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

func CatchException(handle func(e interface{})) {
	if err := recover(); err != nil {
		e := printStackTrace(err)
		handle(e)
	}
}

// 打印堆栈信息
func printStackTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}

func GetErrorStack(err error, preStr string) string {
	buf := new(bytes.Buffer)
	var errPrefix = "堆栈信息写入buf失败"
	_, errPrint := fmt.Fprintf(buf, "%s\n", preStr)
	if errPrint != nil {
		fmt.Println(errPrint)
		return errPrefix + errPrint.Error()
	}
	_, errPrint = fmt.Fprintf(buf, "Cause by: %v\n", err)
	if errPrint != nil {
		fmt.Println(errPrint)
		return errPrefix + errPrint.Error()
	}
	var fileName string
	var subFile []string
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if !ok {
			break
		}
		fileName = file[strings.LastIndex(file, SLASH) + 1:strings.LastIndex(file, DOT)]
		subFile = strings.Split(f.Name(), ".")
		_, errPrint = fmt.Fprintf(buf, fmt.Sprintf("\tat %s:%d\n", subFile[0] + DOT + fileName + DOT+ subFile[1], line))
		if errPrint != nil {
			fmt.Println(errPrint)
			return errPrefix + errPrint.Error()
		}
	}

	return buf.String()
}

