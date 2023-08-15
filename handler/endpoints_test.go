package handler

import (
	//"context"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/bwmarrin/snowflake"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEndpoints_UserRegistration(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("", "", nil)
	ctx := e.NewContext(req, rec)
	type args struct {
		ctx        echo.Context
	}

	type fields struct {
		Repo repository.RepositoryInterface
		Sf   *snowflake.Node
	}

	tests := []struct {
		name           string
		args           args
		// expectedResult generated.RegisterResponse
		expectedErr    error
		mockFn         func(a args, expectedErr error) fields
	}{
		{
			name: "ErrorJSONDecoder",
			args: args{
				ctx:        ctx,
			},
			expectedErr:    nil,
			mockFn: func(a args, expectedErr error) fields {
				return fields{
					Repo: repository.NewMockRepositoryInterface(&gomock.Controller{}),
					Sf:   nil,
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := test.mockFn(test.args, test.expectedErr)

			srv := &Server{
				Repository: app.Repo,
				Snowflake:  app.Sf,
			}

			err := srv.UserRegistration(test.args.ctx)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
