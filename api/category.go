package api

import (
	"net/http"

	"github.com/d1vshar/splitgo/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (Api *Api) GetCategoriesHandler(c echo.Context) error {
	res, err := Api.App.DB().GetCategories()

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessResponse(c, http.StatusOK, res)
}

func (Api *Api) CreateCategoryHandler(c echo.Context) error {
	t := new(dto.NewCategoryBody)
	if err := c.Bind(t); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	if new_category, err := Api.App.DB().CreateCategory(t.Name); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	} else {
		return SuccessResponse(c, http.StatusOK, new_category)
	}
}

func (Api *Api) UpdateCategoryHandler(c echo.Context) error {
	id := c.Param("id")

	categoryID, err := uuid.Parse(id)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	t := new(dto.NewCategoryBody)
	if err := c.Bind(t); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	if category, err := Api.App.DB().UpdateCategory(categoryID, t.Name); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	} else {
		return SuccessResponse(c, http.StatusOK, category)
	}
}

func (Api *Api) DeleteCategoryHandler(c echo.Context) error {
	id := c.Param("id")

	categoryID, err := uuid.Parse(id)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err)
	}

	if err := Api.App.DB().DeleteCategory(categoryID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessResponse(c, http.StatusOK, nil)
}

func (Api *Api) RegisterCategoryRoutes(g *echo.Group) {
	var tg = g.Group("/category")

	tg.GET("", Api.GetCategoriesHandler)
	tg.POST("", Api.CreateCategoryHandler)
	tg.PUT("/:id", Api.UpdateCategoryHandler)
	tg.DELETE("/:id", Api.DeleteCategoryHandler)
}
