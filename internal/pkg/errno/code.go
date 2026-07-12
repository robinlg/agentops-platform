package errno

import (
	"net/http"

	"github.com/robinlg/onexlib/pkg/errorsx"
)

var (
	// OK 代表请求成功
	OK = &errorsx.ErrorX{Code: http.StatusOK, Message: ""}

	// ErrInternal 表示所有未知的服务器端错误
	ErrInternal = errorsx.ErrInternal

	// ErrPageNotFound 表示页面未找到
	ErrPageNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PageNotFound", Message: "Page not found."}
)
