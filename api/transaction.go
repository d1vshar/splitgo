package api

import (
	"net/http"

	"github.com/d1vshar/splitgo/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (Api *Api) GetTransactionsHandler(c echo.Context) error {
	trxs, err := Api.App.DB().GetTransactionsForUser("0ec02773-f863-400f-9c45-2a30b1e547b7")

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessResponse(c, http.StatusOK, trxs)
}

func (Api *Api) CreateTransactionHandler(c echo.Context) error {
	t := new(dto.TransasctionPayload)
	if err := c.Bind(t); err != nil {
		return err
	}

	res, err := Api.App.DB().CreateNewTransaction("0ec02773-f863-400f-9c45-2a30b1e547b7", t)

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessResponse(c, http.StatusOK, res)
}

func (Api *Api) UpdateTransactionHandler(c echo.Context) error {
	id := c.Param("id")
	trx_id, err := uuid.Parse(id)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	t := new(dto.TransasctionPayload)
	if err := c.Bind(t); err != nil {
		return err
	}

	res, err := Api.App.DB().UpdateTransaction(trx_id, "0ec02773-f863-400f-9c45-2a30b1e547b7", t)

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessResponse(c, http.StatusOK, res)
}

func (Api *Api) DeleteTransactionHandler(c echo.Context) error {
	id := c.Param("id")
	trx_id, err := uuid.Parse(id)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	delete_err := Api.App.DB().DeleteTransaction(trx_id, "0ec02773-f863-400f-9c45-2a30b1e547b7")

	if delete_err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, delete_err)
	}

	return SuccessResponse(c, http.StatusOK, nil)
}

func (Api *Api) RegisterTransactionRoutes(g *echo.Group) {
	var tg = g.Group("/transaction")

	tg.GET("", Api.GetTransactionsHandler)
	tg.POST("", Api.CreateTransactionHandler)
	tg.PUT("/:id", Api.UpdateTransactionHandler)
	tg.DELETE("/:id", Api.DeleteTransactionHandler)
}
