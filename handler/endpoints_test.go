package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	server *Server

	mockRepository *repository.MockRepositoryInterface
)

func provideTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository = repository.NewMockRepositoryInterface(ctrl)
	server = NewServer(NewServerOptions{
		Repository: mockRepository,
		SecretKey:  "sawitpro",
	})

	return func() {}
}

func TestLogin(t *testing.T) {
	t.Run("TestLogin", func(t *testing.T) {
		Convey("TestLogin", t, func(c C) {
			type (
				args struct {
					payload string
				}
			)

			testCases := []struct {
				testID         int
				testDesc       string
				args           args
				mockFunc       func()
				wantStatusCode int
				wantResp       generated.LoginResponse
			}{
				{
					testID:   1,
					testDesc: "Failed - error Bind",
					args: args{
						payload: `{"phone":"+6280989444","password"s:"password-mock"}`,
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusBadRequest,
				},
				{
					testID:   2,
					testDesc: "Failed - error GetUserByPhone",
					args: args{
						payload: `{"phone":"+6280989444","password":"password-mock"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{}, fmt.Errorf("error"))
					},
					wantStatusCode: http.StatusBadRequest,
				},
				{
					testID:   3,
					testDesc: "Failed - error ComparePassword",
					args: args{
						payload: `{"phone":"+6280989444","password":"password1!A"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{
							ID:       1,
							Password: "$2a$04$eMb1vD6rv6hXe/PKA2Wzj.b1dO0oW2PTYQzA5ez8Rm3sGrD6ULrKd2",
						}, nil)
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.LoginResponse{},
				},
				{
					testID:   4,
					testDesc: "Failed - error IncreaseLoginCount",
					args: args{
						payload: `{"phone":"+6280989444","password":"password1!A"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{
							ID:       1,
							Password: "$2a$04$eMb1vD6rv6hXe/PKA2Wzj.b1dO0oW2PTYQzA5ez8Rm3GrD6ULrKd2",
						}, nil)
						mockRepository.EXPECT().IncreaseLoginCount(gomock.Any(), int64(1)).Return(fmt.Errorf("error"))
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.LoginResponse{},
				},
				{
					testID:   5,
					testDesc: "Success",
					args: args{
						payload: `{"phone":"+6280989444","password":"password1!A"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{
							ID:       1,
							Password: "$2a$04$eMb1vD6rv6hXe/PKA2Wzj.b1dO0oW2PTYQzA5ez8Rm3GrD6ULrKd2",
						}, nil)
						mockRepository.EXPECT().IncreaseLoginCount(gomock.Any(), int64(1)).Return(nil)
					},
					wantStatusCode: http.StatusOK,
					wantResp: generated.LoginResponse{
						Id: 1,
					},
				},
			}

			for _, tc := range testCases {
				testDep := provideTest(t)
				defer testDep()

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					tc.mockFunc()

					method := echo.POST
					path := "/login"

					e := echo.New()
					req := httptest.NewRequest(method, path, bytes.NewReader([]byte(tc.args.payload)))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rr := httptest.NewRecorder()
					c := e.NewContext(req, rr)
					_ = server.Login(c)

					// assert
					var resp generated.LoginResponse
					_ = json.Unmarshal(rr.Body.Bytes(), &resp)
					So(resp.Id, ShouldEqual, tc.wantResp.Id)
					So(rr.Code, ShouldEqual, tc.wantStatusCode)
				})
			}
		})
	})
}

func TestUserRegister(t *testing.T) {
	t.Run("TestUserRegister", func(t *testing.T) {
		Convey("TestUserRegister", t, func(c C) {
			type (
				args struct {
					payload string
				}
			)

			testCases := []struct {
				testID         int
				testDesc       string
				args           args
				mockFunc       func()
				wantStatusCode int
				wantResp       generated.RegisterResponse
			}{
				{
					testID:   1,
					testDesc: "Failed - error Bind",
					args: args{
						payload: `{"phone":"+6280989444","name":"albert einstein",2"password"s:"password-mock"}`,
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusBadRequest,
				},
				{
					testID:   2,
					testDesc: "Failed - error Validation",
					args: args{
						payload: `{"phone":"+6280989444","name":"albert einstein","password":"Password"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{}, nil)
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.RegisterResponse{},
				},
				{
					testID:   3,
					testDesc: "Failed - phone number duplicate",
					args: args{
						payload: `{"phone":"+6280989444","name":"albert einstein","password":"Password1!"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{}, nil)
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.RegisterResponse{},
				},
				{
					testID:   4,
					testDesc: "Failed - error Createuser",
					args: args{
						payload: `{"phone":"+6280989444","name":"albert einstein","password":"Password1!"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{}, fmt.Errorf("error"))
						mockRepository.EXPECT().Createuser(gomock.Any(), gomock.Any()).Return(repository.User{}, fmt.Errorf("error"))
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.RegisterResponse{},
				},
				{
					testID:   5,
					testDesc: "Success",
					args: args{
						payload: `{"phone":"+6280989444","name":"albert einstein","password":"Password1!"}`,
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{}, fmt.Errorf("error"))
						mockRepository.EXPECT().Createuser(gomock.Any(), gomock.Any()).Return(repository.User{
							ID: 1,
						}, nil)
					},
					wantStatusCode: http.StatusOK,
					wantResp: generated.RegisterResponse{
						Id: 1,
					},
				},
			}

			for _, tc := range testCases {
				testDep := provideTest(t)
				defer testDep()

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					tc.mockFunc()

					method := echo.POST
					path := "/users/register"

					e := echo.New()
					req := httptest.NewRequest(method, path, bytes.NewReader([]byte(tc.args.payload)))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					rr := httptest.NewRecorder()
					c := e.NewContext(req, rr)
					_ = server.UserRegister(c)

					// assert
					var resp generated.RegisterResponse
					_ = json.Unmarshal(rr.Body.Bytes(), &resp)
					So(resp.Id, ShouldEqual, tc.wantResp.Id)
					So(rr.Code, ShouldEqual, tc.wantStatusCode)
				})
			}
		})
	})
}

func TestGetUser(t *testing.T) {
	t.Run("TestGetUser", func(t *testing.T) {
		Convey("TestGetUser", t, func(c C) {
			type (
				args struct {
					authorization string
				}
			)

			testCases := []struct {
				testID         int
				testDesc       string
				args           args
				mockFunc       func()
				wantStatusCode int
				wantResp       generated.UserResponse
			}{
				{
					testID:   1,
					testDesc: "Failed - error ValidateJWT",
					args: args{
						authorization: "Bearer 1eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   2,
					testDesc: "Failed - error ValidateJWT - invalid auth",
					args: args{
						authorization: "BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   3,
					testDesc: "Failed - error ValidateJWT - invalid auth",
					args: args{
						authorization: "Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   4,
					testDesc: "Failed - error ValidateJWT - invalid auth not integer",
					args: args{
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMWIifQ.4fWpYB6cdmowGR1g2dTiCPsYQ_8X-q_2oftgIsEy8cE",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   5,
					testDesc: "Failed - error GetUserByID",
					args: args{
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByID(gomock.Any(), int64(17)).Return(repository.User{}, fmt.Errorf("error"))
					},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   6,
					testDesc: "Success",
					args: args{
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc: func() {
						mockRepository.EXPECT().GetUserByID(gomock.Any(), int64(17)).Return(repository.User{
							Phone: "+62812922222",
							Name:  "mr mozart1",
						}, nil)
					},
					wantStatusCode: http.StatusOK,
					wantResp: generated.UserResponse{
						Phone: "+62812922222",
						Name:  "mr mozart1",
					},
				},
			}

			for _, tc := range testCases {
				testDep := provideTest(t)
				defer testDep()

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					tc.mockFunc()

					method := echo.GET
					path := "/users"

					e := echo.New()
					req := httptest.NewRequest(method, path, bytes.NewReader([]byte("")))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					req.Header.Set(echo.HeaderAuthorization, tc.args.authorization)
					rr := httptest.NewRecorder()
					c := e.NewContext(req, rr)
					_ = server.GetUser(c)

					// assert
					var resp generated.UserResponse
					_ = json.Unmarshal(rr.Body.Bytes(), &resp)
					So(resp.Phone, ShouldEqual, tc.wantResp.Phone)
					So(resp.Name, ShouldEqual, tc.wantResp.Name)
					So(rr.Code, ShouldEqual, tc.wantStatusCode)
				})
			}
		})
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("TestUpdateUser", func(t *testing.T) {
		Convey("TestUpdateUser", t, func(c C) {
			type (
				args struct {
					authorization string
					payload       string
				}
			)

			testCases := []struct {
				testID         int
				testDesc       string
				args           args
				mockFunc       func()
				wantStatusCode int
				wantResp       generated.UserResponse
			}{
				{
					testID:   1,
					testDesc: "Failed - error ValidateJWT",
					args: args{
						payload:       `{"phone":"+6280989444","name":"halo halo"}`,
						authorization: "Bearer 1eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   2,
					testDesc: "Failed - error ValidateJWT - not integer",
					args: args{
						payload:       `{"phone":"+6280989444","name":"halo halo"}`,
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMWIifQ.4fWpYB6cdmowGR1g2dTiCPsYQ_8X-q_2oftgIsEy8cE",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusForbidden,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   3,
					testDesc: "Failed - error Bind",
					args: args{
						payload:       `{"phone":"+6280989444","name":"halo halo"1}`,
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc:       func() {},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   4,
					testDesc: "Failed - error validate",
					args: args{
						payload:       `{"phone":"+628098","name":"halo halo"}`,
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc: func() {
						mockRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(repository.User{}, nil)
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   5,
					testDesc: "Failed - error UpdateUser",
					args: args{
						payload:       `{"phone":"+6280989444","name":"halo halo"}`,
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc: func() {
						mockRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(repository.User{}, fmt.Errorf("error"))
					},
					wantStatusCode: http.StatusBadRequest,
					wantResp:       generated.UserResponse{},
				},
				{
					testID:   6,
					testDesc: "Success",
					args: args{
						payload:       `{"phone":"+6280989444","name":"halo halo"}`,
						authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTcifQ.McgHf0xflNEkFStawm8KgdjyPHQ3he-jI_h8vqBqTZQ",
					},
					mockFunc: func() {
						mockRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(repository.User{
							Phone: "+6280989444",
							Name:  "halo halo",
						}, nil)
					},
					wantStatusCode: http.StatusOK,
					wantResp: generated.UserResponse{
						Phone: "+6280989444",
						Name:  "halo halo",
					},
				},
			}

			for _, tc := range testCases {
				testDep := provideTest(t)
				defer testDep()

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					tc.mockFunc()

					method := echo.PATCH
					path := "/users"

					e := echo.New()
					req := httptest.NewRequest(method, path, bytes.NewReader([]byte(tc.args.payload)))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					req.Header.Set(echo.HeaderAuthorization, tc.args.authorization)
					rr := httptest.NewRecorder()
					c := e.NewContext(req, rr)
					_ = server.UpdateUser(c)

					// assert
					var resp generated.UserResponse
					_ = json.Unmarshal(rr.Body.Bytes(), &resp)
					So(resp.Phone, ShouldEqual, tc.wantResp.Phone)
					So(resp.Name, ShouldEqual, tc.wantResp.Name)
					So(rr.Code, ShouldEqual, tc.wantStatusCode)
				})
			}
		})
	})
}
