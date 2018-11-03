package main

import (
	"net/http"

	"github.com/goadesign/goa"
)

type goaContextWithNoContent interface {
	NoContent() error
}

type goaContextWithBadRequest interface {
	BadRequest() error
}

type goaContextWithBadRequestError interface {
	BadRequest(error) error
}

type goaContextWithNotFoundError interface {
	NotFound(error) error
}

type goaContextWithNotFound interface {
	NotFound() error
}

type goaContextWithRunningContainer interface {
	RunningContainer() error
}

type goaContextWithConflict interface {
	Conflict(error) error
}

type goaContextWithRequestEntityTooLarge interface {
	RequestEntityTooLarge() error
}

type goaContextWithInternalServerError interface {
	InternalServerError(error) error
}

type goaContextHandlerFunctionType func(int, error) error

type errorHandler struct {
	Context interface{}
	method  goaContextHandlerFunctionType
}

func newErrorHandler(ctx interface{}) *errorHandler {
	return &errorHandler{
		Context: ctx,
	}
}

func (eh *errorHandler) handleNoContent() *errorHandler {
	ctx := eh.Context.(goaContextWithNoContent)

	if ctx == nil {
		return nil
	}

	return eh.handle(http.StatusNoContent, ctx.NoContent)
}

func (eh *errorHandler) handleBadRequest() *errorHandler {
	ctx := eh.Context.(goaContextWithBadRequest)

	if ctx == nil {
		return nil
	}

	return eh.handle(http.StatusBadRequest, ctx.BadRequest)
}

func (eh *errorHandler) handleBadRequestWithError() *errorHandler {
	ctx := eh.Context.(goaContextWithBadRequestError)

	if ctx == nil {
		return nil
	}

	return eh.handleWithError(http.StatusBadRequest, ctx.BadRequest)
}

func (eh *errorHandler) handleNotFound() *errorHandler {
	ctx := eh.Context.(goaContextWithNotFound)

	if ctx == nil {
		return nil
	}

	return eh.handle(http.StatusNotFound, ctx.NotFound)
}

func (eh *errorHandler) handleNotFoundWithError() *errorHandler {
	ctx := eh.Context.(goaContextWithNotFoundError)

	if ctx == nil {
		return nil
	}

	return eh.handleWithError(http.StatusNotFound, ctx.NotFound)
}

func (eh *errorHandler) handleRunningContainer() *errorHandler {
	ctx := eh.Context.(goaContextWithRunningContainer)

	if ctx == nil {
		return nil
	}

	return eh.handle(409, ctx.RunningContainer)
}

func (eh *errorHandler) handleConflict() *errorHandler {
	ctx := eh.Context.(goaContextWithConflict)

	if ctx == nil {
		return nil
	}

	return eh.handleWithError(409, ctx.Conflict)
}

func (eh *errorHandler) handleRequestEntityTooLarge() *errorHandler {
	ctx := eh.Context.(goaContextWithRequestEntityTooLarge)

	if ctx == nil {
		return nil
	}

	return eh.handle(http.StatusRequestEntityTooLarge, ctx.RequestEntityTooLarge)
}

func (eh *errorHandler) handleInternalServerError() *errorHandler {
	ctx, ok := eh.Context.(goaContextWithInternalServerError)

	if !ok {
		return nil
	}

	return eh.handleWithError(http.StatusInternalServerError, func(err error) error {
		return ctx.InternalServerError(goa.ErrInternal(err))
	})
}

func (eh *errorHandler) handleWithError(status int, f func(error) error) *errorHandler {
	method := eh.method

	eh.method = func(statusCode int, err error) error {
		if statusCode == status {
			return f(err)
		}

		if method == nil {
			return nil
		}

		return method(statusCode, err)
	}

	return eh
}

func (eh *errorHandler) handle(status int, f func() error) *errorHandler {
	method := eh.method

	eh.method = func(statusCode int, err error) error {
		if statusCode == status {
			return f()
		}

		if method == nil {
			return nil
		}

		return method(statusCode, err)
	}

	return eh
}

func (eh *errorHandler) Call(statusCode int, err error) error {
	if eh.method == nil {
		return nil
	}
	return eh.method(statusCode, err)
}

func isError(status int) bool {
	return status/100 > 2
}
