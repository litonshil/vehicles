package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"vehicles/internal/domain"
	"vehicles/types"
)

type UserController struct {
	baseCtx context.Context
	useruc  domain.UserUseCase
}

func NewUserController(
	baseCtx context.Context,
	useruc domain.UserUseCase,
) *UserController {
	return &UserController{
		baseCtx: baseCtx,
		useruc:  useruc,
	}
}
func (cntlr *UserController) CreateUser(c echo.Context) error {
	var user types.UserReq

	// Bind JSON payload to user struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Call service layer
	err := cntlr.useruc.CreateUser(context.TODO(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func (cntlr *UserController) GetUsers(c echo.Context) error {
	// Call service layer
	userID := c.QueryParam("id")
	filter := domain.UserFilter{
		ID: userID,
	}
	users, err := cntlr.useruc.GetUsers(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch users",
		})
	}

	// Respond with results
	return c.JSON(http.StatusOK, users)
}
