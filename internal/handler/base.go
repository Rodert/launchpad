package handler

import (
	"launchpad/constant/serror"
	base "launchpad/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
	MethodHEAD   = "HEAD"
)

func render(c *gin.Context, v any, err error) {
	resp := base.Response{Data: v, Code: serror.Success, Msg: serror.SuccessMsg}
	if err != nil {
		if tem, ok := err.(*serror.Error); ok {
			resp = base.Response{Data: v, Code: tem.Code, Msg: tem.Error()}
		} else {
			resp = base.Response{Data: v, Code: serror.Unknown, Msg: serror.UnknownMsg}
		}
	}
	c.JSON(http.StatusOK, resp)
}

type (
	TRPathParamHandlerFunc[T any, R any] func(ctx *gin.Context, t *T) (R, error)
)

func TRPathParamHandler[T any, R any](
	handler TRPathParamHandlerFunc[T, R],
) gin.HandlerFunc {
	return handlerWithContext[T, R](handler)
}

func handlerWithContext[T any, R any](
	handlerFunc any,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := new(T)
		if err := c.ShouldBind(&t); err != nil {
			render(c, nil, serror.New(serror.ParamsErr))
			return
		}
		switch handler := handlerFunc.(type) {
		case TRPathParamHandlerFunc[T, R]:
			v, e := handler(c, t)
			render(c, v, e)
		}

	}
}

// 原样返回，不包装
func TRPathParamNudeHandler[T any, R any](
	handler TRPathParamHandlerFunc[T, R],
) gin.HandlerFunc {
	return handlerWithNudeContext[T, R](handler)
}

func handlerWithNudeContext[T any, R any](
	handlerFunc any,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := new(T)
		if err := c.ShouldBind(&t); err != nil {
			renderNude(c, nil, serror.New(serror.ParamsErr))
			return
		}
		switch handler := handlerFunc.(type) {
		case TRPathParamHandlerFunc[T, R]:
			v, e := handler(c, t)
			renderNude(c, v, e)
		}

	}
}

func renderNude(c *gin.Context, v any, err error) {
	resp := v
	c.JSON(http.StatusOK, resp)
}
