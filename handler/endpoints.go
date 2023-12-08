package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

// This function handles user registration.
// (POST /users/register)
func (s *Server) UserRegister(ctx echo.Context) error {
	_ctx := ctx.Request().Context()
	var payload generated.RegisterRequest
	err := ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	input := repository.RegisterUser{
		Phone:    strings.TrimSpace(payload.Phone),
		Name:     strings.TrimSpace(payload.Name),
		Password: strings.TrimSpace(payload.Password),
	}

	if errs := input.Validate(); len(errs) > 0 {
		var resp generated.ErrorResponse
		resp.Message = strings.Join(errs, ", ")
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	input.Password, _ = HashPassword(input.Password)
	_, err = s.Repository.GetUserByPhone(_ctx, payload.Phone)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp, err := s.Repository.Createuser(_ctx, input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.RegisterResponse{
		Id: resp.ID,
	})
}

// This function returns user information.
// (GET /users)
func (s *Server) GetUser(ctx echo.Context) error {
	_ctx := ctx.Request().Context()
	auth := ctx.Request().Header.Get("Authorization")
	token, err := s.ValidateJWT(auth)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	userID, err := strconv.Atoi(s.GetJWTClaims(token, "user_id"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	resp, err := s.Repository.GetUserByID(_ctx, int64(userID))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.UserResponse{
		Phone: resp.Phone,
		Name:  resp.Name,
	})
}

// This function handles user data updation.
// (GET /users)
func (s *Server) UpdateUser(ctx echo.Context) error {
	_ctx := ctx.Request().Context()
	auth := ctx.Request().Header.Get("Authorization")
	token, err := s.ValidateJWT(auth)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	userID, err := strconv.Atoi(s.GetJWTClaims(token, "user_id"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	var payload generated.UpdateUserRequest
	err = ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	input := repository.User{
		ID:    int64(userID),
		Name:  strings.TrimSpace(payload.Name),
		Phone: strings.TrimSpace(payload.Phone),
	}

	if errs := input.Validate(); len(errs) > 0 {
		var resp generated.ErrorResponse
		resp.Message = strings.Join(errs, ", ")
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	resp, err := s.Repository.UpdateUser(_ctx, input)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.UserResponse{
		Phone: resp.Phone,
		Name:  resp.Name,
	})
}

// This function handles user login.
// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	_ctx := ctx.Request().Context()
	var payload generated.LoginRequest
	err := ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	user, err := s.Repository.GetUserByPhone(_ctx, payload.Phone)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = ComparePassword(payload.Password, user.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "incorrect password or phone number",
		})
	}

	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = s.Repository.IncreaseLoginCount(_ctx, user.ID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Id:    user.ID,
		Token: token,
	})
}
