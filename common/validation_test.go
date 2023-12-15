package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateName(t *testing.T) {
	type (
		args struct {
			payload string
		}
	)

	testCases := []struct {
		testID   int
		testDesc string
		args     args
		wantErr  bool
	}{
		{
			testID:   1,
			testDesc: "Failed - name too short",
			args: args{
				payload: `ab`,
			},
			wantErr: true,
		},
		{
			testID:   2,
			testDesc: "Failed - name too long",
			args: args{
				payload: `integer malesuada nunc vel risus commodo viverra maecenas accumsan lacus vel facilisis volutpat est velit egestas dui id ornare arcu`,
			},
			wantErr: true,
		},
		{
			testID:   5,
			testDesc: "Success",
			args: args{
				payload: `integer malesuada`,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			resp := ValidateName(tc.args.payload)
			assert.Equal(t, len(resp) > 0, tc.wantErr)
		})
	}
}

func TestValidatePhone(t *testing.T) {
	type (
		args struct {
			payload string
		}
	)

	testCases := []struct {
		testID   int
		testDesc string
		args     args
		wantErr  bool
	}{
		{
			testID:   1,
			testDesc: "Failed - error prefix number",
			args: args{
				payload: `+61809A89444`,
			},
			wantErr: true,
		},
		{
			testID:   2,
			testDesc: "Failed - contains non number",
			args: args{
				payload: `+6280989444A`,
			},
			wantErr: true,
		},
		{
			testID:   3,
			testDesc: "Failed - phone too short",
			args: args{
				payload: `+628098`,
			},
			wantErr: true,
		},
		{
			testID:   4,
			testDesc: "Failed - phone too long",
			args: args{
				payload: `+628098944412121`,
			},
			wantErr: true,
		},
		{
			testID:   5,
			testDesc: "Success",
			args: args{
				payload: `+6280989444`,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			resp := ValidatePhone(tc.args.payload)
			assert.Equal(t, len(resp) > 0, tc.wantErr)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	type (
		args struct {
			payload string
		}
	)

	testCases := []struct {
		testID   int
		testDesc string
		args     args
		wantErr  bool
	}{
		{
			testID:   1,
			testDesc: "Failed - password too short",
			args: args{
				payload: `Pas`,
			},
			wantErr: true,
		},
		{
			testID:   2,
			testDesc: "Failed - password too long",
			args: args{
				payload: `integer malesuada nunc vel risus commodo viverra maecenas accumsan lacus vel facilisis volutpat est velit egestas dui id ornare arcu`,
			},
			wantErr: true,
		},
		{
			testID:   3,
			testDesc: "Failed - not contains capital",
			args: args{
				payload: `password1!`,
			},
			wantErr: true,
		},
		{
			testID:   4,
			testDesc: "Success",
			args: args{
				payload: `Password1!`,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDesc, func(t *testing.T) {
			resp := ValidatePassword(tc.args.payload)
			assert.Equal(t, len(resp) > 0, tc.wantErr)
		})
	}
}
