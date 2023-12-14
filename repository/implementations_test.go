package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateuser(t *testing.T) {
	t.Run("TestLogin", func(t *testing.T) {
		Convey("TestLogin", t, func(c C) {
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
					testID:   5,
					testDesc: "Success",
					args: args{
						payload: RegisterUser{
							Name:     "mock-name",
							Phone:    "mock-phone",
							Password: "mock-password",
						}, //`{"phone":"+6280989444","password":"password1!A"}`,
					},
					mockFunc: func(mockSQL sqlmock.Sqlmock) {

						// VALUES ($1, $2, $3)
						// RETURNING id
						mockSQL.ExpectBegin()
						mockSQL.ExpectPrepare(`INSERT INTO users (.+)`)
						mockSQL.ExpectQuery("INSERT INTO users (.+)").
							WithArgs("mock-phone", "mock-name", "mock-password").
							WillReturnRows(
								sqlmock.NewRows([]string{"id"}).
									AddRow(1))
						mockSQL.ExpectCommit()
						// mo// mockRepository.EXPECT().GetUserByPhone(gomock.Any(), "+6280989444").Return(repository.User{
						// // 	ID:       1,
						// 	Password: "$2a$04$eMb1vD6rv6hXe/PKA2Wzj.b1dO0oW2PTYQzA5ez8Rm3GrD6ULrKd2",
						// }, nil)
						// mockRepository.EXPECT().IncreaseLoginCount(gomock.Any(), int64(1)).Return(nil)
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

					So(output, ShouldEqual, tc.wantResp)
					So(err != nil, ShouldEqual, tc.wantErr)
				})
			}
		})
	})
}
