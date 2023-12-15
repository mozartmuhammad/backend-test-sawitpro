package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

// POST API create new user data.
// http://localhost:1323/users/register
func (s *Server) UserRegister(c echo.Context) error {
	ctx := c.Request().Context()
	var payload generated.RegisterRequest
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	input := repository.RegisterUser{
		Phone:    strings.TrimSpace(payload.Phone),
		Name:     strings.TrimSpace(payload.Name),
		Password: strings.TrimSpace(payload.Password),
	}

	// validate input with given rules.
	if errs := input.Validate(); len(errs) > 0 {
		var resp generated.ErrorResponse
		resp.Message = strings.Join(errs, ", ")
		return c.JSON(http.StatusBadRequest, resp)
	}

	// check whether user phone exist.
	_, err = s.Repository.GetUserByPhone(ctx, payload.Phone)
	if err == nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phone number already exist",
		})
	}

	// hash and salt user password.
	input.Password, _ = HashPassword(input.Password)

	// create user data.
	resp, err := s.Repository.Createuser(ctx, input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.RegisterResponse{
		Id: resp.ID,
	})
}

// GET API which return the name and phone number of the cats based on authorization.
// http://localhost:1323/users
func (s *Server) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	auth := c.Request().Header.Get("Authorization")

	// validate authorization.
	token, err := s.ValidateJWT(auth)
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	// get user_id from token.
	userID, err := strconv.Atoi(s.GetJWTClaims(token, "user_id"))
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	// get user data by id.
	resp, err := s.Repository.GetUserByID(ctx, int64(userID))
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.UserResponse{
		Phone: resp.Phone,
		Name:  resp.Name,
	})
}

// PATCH API responsible to update user data.
// http://localhost:1323/users
func (s *Server) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	auth := c.Request().Header.Get("Authorization")

	// validate authorization.
	token, err := s.ValidateJWT(auth)
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	// get user_id from token.
	userID, err := strconv.Atoi(s.GetJWTClaims(token, "user_id"))
	if err != nil {
		return c.JSON(http.StatusForbidden, generated.ErrorResponse{Message: err.Error()})
	}

	var payload generated.UpdateUserRequest
	err = c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	input := repository.User{
		ID:    int64(userID),
		Name:  strings.TrimSpace(payload.Name),
		Phone: strings.TrimSpace(payload.Phone),
	}

	// validate input with given rules.
	if errs := input.Validate(); len(errs) > 0 {
		var resp generated.ErrorResponse
		resp.Message = strings.Join(errs, ", ")
		return c.JSON(http.StatusBadRequest, resp)
	}

	// process update user data.
	resp, err := s.Repository.UpdateUser(ctx, input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.UserResponse{
		Phone: resp.Phone,
		Name:  resp.Name,
	})
}

// POST API responsible to log in user.
// http://localhost:1323/login
func (s *Server) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload generated.LoginRequest
	err := c.Bind(&payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	// check whether phone number exist.
	user, err := s.Repository.GetUserByPhone(ctx, payload.Phone)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// compare user password.
	err = ComparePassword(payload.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "incorrect password or phone number",
		})
	}

	// generate JWT.
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// increate user login count.
	err = s.Repository.IncreaseLoginCount(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, generated.LoginResponse{
		Id:    user.ID,
		Token: token,
	})
}
