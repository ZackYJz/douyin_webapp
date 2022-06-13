package app

import (
	"go_webapp/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Ctx *gin.Context
}

type ResponseData struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

type Pagination struct {
	Page       int         `json:"page"`    //页码 （当前页数）
	TotalRows  int         `json:"records"` //总记录行数
	TotalPages int         `json:"total"`   //总页数
	Data       interface{} `json:"rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ResponseSuccess(data interface{}) {
	r.Ctx.JSON(http.StatusOK, &ResponseData{
		Status: errcode.Success.Code(),
		Msg:    errcode.Success.Msg(),
		Data:   data,
	})
}

func (r *Response) ResponseList(data interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, &ResponseData{
		Status: errcode.Success.Code(),
		Msg:    errcode.Success.Msg(),
		Data: &Pagination{
			Page:       GetPage(r.Ctx),
			TotalPages: totalRows/GetPageSize(r.Ctx) + 1,
			TotalRows:  totalRows,
			Data:       data,
		},
	})
}

func (r *Response) ResponseError(err *errcode.Error) {

	r.Ctx.JSON(http.StatusOK, &ResponseData{
		Status: err.Code(),
		Msg:    err.Msg(),
		Data:   nil,
	})
}

func (r *Response) ResponseErrorString(errStr string) {

	r.Ctx.JSON(http.StatusOK, &ResponseData{
		Status: errcode.ServerError.Code(),
		Msg:    errStr,
		Data:   nil,
	})
}
