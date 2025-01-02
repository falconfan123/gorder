package common

import (
	"github.com/falconfan123/gorder/common/tracing"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct{}

// 小写，不导出他
type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
	if err != nil {
		base.Error(c, err)
	} else {
		base.Success(c, data)
	}
}

func (base *BaseResponse) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Errno:   0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}

func (base *BaseResponse) Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
