package service

import (
	"context"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/entity/user"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/mongodb"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/jwt"
	authpb "github.com/mhdiiilham/oauth2-auth-server-implementation/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestRegisterNewUser(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name    string
		code    codes.Code
		payload *authpb.ResgisterRequest
	}{
		{
			name: "Success registering new user",
			code: codes.OK,
			payload: &authpb.ResgisterRequest{
				Fullname: "User Testing",
				Email:    "Testing@mail.com",
				Password: "PasswordTest.com",
			},
		},
		{
			name:    "Failed registering new user",
			code:    codes.InvalidArgument,
			payload: &authpb.ResgisterRequest{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, coll, mongoErr := mongodb.NewMongoDBConnection(os.Getenv("MONGO_DB_USER"), os.Getenv("MONGO_DB_PASS"), os.Getenv("MONGO_DB"), os.Getenv("MONGO_DB_COLLECTION"))
			if mongoErr != nil {
				panic(mongoErr)
			}

			tokenService := jwt.NewJWTService(string(os.Getenv("JWT_SECRET")), string(os.Getenv("APP_NAME")))
			userRepo := user.NewMongoDBRepository(coll)
			userManager := user.NewManager(userRepo)
			authServer := NewService(os.Getenv("SERVER_NETWORK"), os.Getenv("SERVER_ADDRESS"), userManager, tokenService)

			res, err := authServer.RegisterService(context.TODO(), &authpb.ResgisterRequest{
				Fullname: tc.payload.GetFullname(),
				Email:    tc.payload.GetEmail(),
				Password: tc.payload.GetPassword(),
			})

			if tc.code == codes.OK {
				assert.Nil(t, err, "Error should be empty")
				assert.NotNil(t, res.GetAccessToken(), "Access token should not be empty")
				assert.Equal(t, res.GetMessage(), "Register success", "Message should be equal to 'Register Success'")
			} else {
				assert.NotNil(t, err, "Error should be not nil")
				assert.Nil(t, res, "Response should be nil")
			}
		})
	}
}
