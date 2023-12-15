package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateuser(t *testing.T) {
	t.Run("TestCreateuser", func(t *testing.T) {
		Convey("TestCreateuser", t, func(c C) {
			mockUser := RegisterUser{
				Name:     "mock-name",
				Phone:    "mock-phone",
				Password: "mock-password",
			}

			type (
				args struct {
					payload RegisterUser
				}
			)

			testCases := []struct {
				testID   int
				testDesc string
				args     args
				mockFunc func(mockSQL sqlmock.Sqlmock)
				wantResp User
				wantErr  bool
			}{
				{
					testID:   1,
					testDesc: "Failed - error begin",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin().WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   2,
					testDesc: "Failed - error prepare",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`).WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   3,
					testDesc: "Failed - error query - unique constraint",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`)
						mockSQL.ExpectQuery("INSERT INTO users (.+)").
							WithArgs("mock-phone", "mock-name", "mock-password").
							WillReturnError(&pq.Error{Code: "23505"})
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   4,
					testDesc: "Failed - error query",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`)
						mockSQL.ExpectQuery("INSERT INTO users (.+)").
							WithArgs("mock-phone", "mock-name", "mock-password").
							WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   5,
					testDesc: "Failed - error commit",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`)
						mockSQL.ExpectQuery("INSERT INTO users (.+)").
							WithArgs("mock-phone", "mock-name", "mock-password").
							WillReturnRows(
								sqlmock.NewRows([]string{"id"}).
									AddRow(1))
						mockSQL.ExpectCommit().WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   6,
					testDesc: "Success",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`)
						mockSQL.ExpectQuery("INSERT INTO users (.+)").
							WithArgs("mock-phone", "mock-name", "mock-password").
							WillReturnRows(
								sqlmock.NewRows([]string{"id"}).
									AddRow(1))
						mockSQL.ExpectCommit()
					},
					wantErr: false,
					wantResp: User{
						ID: 1,
					},
				},
			}

			for _, tc := range testCases {

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					mockDB, mockSQL, _ := sqlmock.New()
					defer mockDB.Close()

					r := Repository{
						Db: mockDB,
					}
					tc.mockFunc(mockSQL)

					output, err := r.Createuser(context.Background(), tc.args.payload)
					// assert
					So(err != nil, ShouldEqual, tc.wantErr)
					So(output, ShouldEqual, tc.wantResp)
				})
			}
		})
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("TestGetUserByID", func(t *testing.T) {
		Convey("TestGetUserByID", t, func(c C) {
			mockTime := time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC)

			type (
				args struct {
					id int64
				}
			)

			testCases := []struct {
				testID   int
				testDesc string
				args     args
				mockFunc func(mockSQL sqlmock.Sqlmock)
				wantResp User
				wantErr  bool
			}{
				{
					testID:   1,
					testDesc: "Failed",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectQuery("SELECT (.+)").
							WithArgs(int64(1)).
							WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   2,
					testDesc: "Success",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectQuery("SELECT (.+)").
							WithArgs(int64(1)).
							WillReturnRows(
								sqlmock.NewRows([]string{"id", "phone", "name", "password", "created_at", "updated_at"}).
									AddRow(int64(1), "mock-phone", "mock-name", "mock-password", mockTime, nil))
					},
					wantErr: false,
					wantResp: User{
						ID:        1,
						Name:      "mock-name",
						Phone:     "mock-phone",
						Password:  "mock-password",
						CreatedAt: mockTime,
						UpdateAt:  nil,
					},
				},
			}

			for _, tc := range testCases {

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					mockDB, mockSQL, _ := sqlmock.New()
					defer mockDB.Close()

					r := Repository{
						Db: mockDB,
					}
					tc.mockFunc(mockSQL)

					output, err := r.GetUserByID(context.Background(), tc.args.id)
					// assert
					So(err != nil, ShouldEqual, tc.wantErr)
					So(output, ShouldEqual, tc.wantResp)
				})
			}
		})
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("TestUpdateUser", func(t *testing.T) {
		Convey("TestUpdateUser", t, func(c C) {
			mockUser := User{
				ID:    1,
				Name:  "mock-name",
				Phone: "mock-phone",
			}

			type (
				args struct {
					payload User
				}
			)

			testCases := []struct {
				testID   int
				testDesc string
				args     args
				mockFunc func(mockSQL sqlmock.Sqlmock)
				wantErr  bool
			}{
				{
					testID:   1,
					testDesc: "Failed - error begin",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin().WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   2,
					testDesc: "Failed - error prepare",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`).WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   3,
					testDesc: "Failed - error exec - unique constraint",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs("mock-phone", "mock-name", 1).
							WillReturnError(&pq.Error{Code: "23505"})
					},
					wantErr: true,
				},
				{
					testID:   4,
					testDesc: "Failed - error exec",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs("mock-phone", "mock-name", 1).
							WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   5,
					testDesc: "Failed - error commit",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs("mock-phone", "mock-name", 1).
							WillReturnResult(sqlmock.NewResult(0, 1))
						mockSQL.ExpectCommit().WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   6,
					testDesc: "Success",
					args: args{
						payload: mockUser,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs("mock-phone", "mock-name", 1).
							WillReturnResult(sqlmock.NewResult(0, 1))
						mockSQL.ExpectCommit()
					},
					wantErr: false,
				},
			}

			for _, tc := range testCases {

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					mockDB, mockSQL, _ := sqlmock.New()
					defer mockDB.Close()

					r := Repository{
						Db: mockDB,
					}
					tc.mockFunc(mockSQL)

					_, err := r.UpdateUser(context.Background(), tc.args.payload)
					// assert
					So(err != nil, ShouldEqual, tc.wantErr)
				})
			}
		})
	})
}

func TestGetUserByPhone(t *testing.T) {
	t.Run("TestGetUserByPhone", func(t *testing.T) {
		Convey("TestGetUserByPhone", t, func(c C) {
			mockTime := time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC)

			type (
				args struct {
					phone string
				}
			)

			testCases := []struct {
				testID   int
				testDesc string
				args     args
				mockFunc func(mockSQL sqlmock.Sqlmock)
				wantResp User
				wantErr  bool
			}{
				{
					testID:   1,
					testDesc: "Failed",
					args: args{
						phone: "mock-phone",
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectQuery("SELECT (.+)").
							WithArgs("mock-phone").
							WillReturnError(fmt.Errorf("error"))
					},
					wantErr:  true,
					wantResp: User{},
				},
				{
					testID:   2,
					testDesc: "Success",
					args: args{
						phone: "mock-phone",
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectQuery("SELECT (.+)").
							WithArgs("mock-phone").
							WillReturnRows(
								sqlmock.NewRows([]string{"id", "phone", "name", "password", "created_at", "updated_at"}).
									AddRow(int64(1), "mock-phone", "mock-name", "mock-password", mockTime, nil))
					},
					wantErr: false,
					wantResp: User{
						ID:        1,
						Name:      "mock-name",
						Phone:     "mock-phone",
						Password:  "mock-password",
						CreatedAt: mockTime,
						UpdateAt:  nil,
					},
				},
			}

			for _, tc := range testCases {

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					mockDB, mockSQL, _ := sqlmock.New()
					defer mockDB.Close()

					r := Repository{
						Db: mockDB,
					}
					tc.mockFunc(mockSQL)

					output, err := r.GetUserByPhone(context.Background(), tc.args.phone)
					// assert
					So(err != nil, ShouldEqual, tc.wantErr)
					So(output, ShouldEqual, tc.wantResp)
				})
			}
		})
	})
}

func TestIncreaseLoginCount(t *testing.T) {
	t.Run("TestIncreaseLoginCount", func(t *testing.T) {
		Convey("TestIncreaseLoginCount", t, func(c C) {

			type (
				args struct {
					id int64
				}
			)

			testCases := []struct {
				testID   int
				testDesc string
				args     args
				mockFunc func(mockSQL sqlmock.Sqlmock)
				wantErr  bool
			}{
				{
					testID:   1,
					testDesc: "Failed - error begin",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin().WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   2,
					testDesc: "Failed - error prepare",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`).WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   3,
					testDesc: "Failed - error exec",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs(1).
							WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   4,
					testDesc: "Failed - error commit",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs(1).
							WillReturnResult(sqlmock.NewResult(0, 1))
						mockSQL.ExpectCommit().WillReturnError(fmt.Errorf("error"))
					},
					wantErr: true,
				},
				{
					testID:   5,
					testDesc: "Success",
					args: args{
						id: 1,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`UPDATE users(.+)`)
						mockSQL.ExpectExec("UPDATE users(.+)").
							WithArgs(1).
							WillReturnResult(sqlmock.NewResult(0, 1))
						mockSQL.ExpectCommit()
					},
					wantErr: false,
				},
			}

			for _, tc := range testCases {

				Convey(fmt.Sprintf("%d : %s", tc.testID, tc.testDesc), func() {
					mockDB, mockSQL, _ := sqlmock.New()
					defer mockDB.Close()

					r := Repository{
						Db: mockDB,
					}
					tc.mockFunc(mockSQL)

					err := r.IncreaseLoginCount(context.Background(), tc.args.id)
					// assert
					So(err != nil, ShouldEqual, tc.wantErr)
				})
			}
		})
	})
}
