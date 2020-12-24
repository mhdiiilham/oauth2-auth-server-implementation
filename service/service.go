package service

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/entity/user"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/jwt"
	pwd "github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/password"
	authpb "github.com/mhdiiilham/oauth2-auth-server-implementation/protos"
)

// Server struct
type Server struct {
	Network string
	Address string
	Manager user.Manager
	Token   jwt.TokenService
}

// NewService function
// to create new gRPC Service
func NewService(network, address string, userManager user.Manager, ts jwt.TokenService) *Server {
	return &Server{
		Network: network,
		Address: address,
		Manager: userManager,
		Token:   ts,
	}
}

/*
|--------------------------------------------------------------------------
| Server Methods
|--------------------------------------------------------------------------
|
| Here is where you can register gRPC Methods.
|
*/

// RegisterService Handler
func (s *Server) RegisterService(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	log.Printf("Registering user with email %s \n", req.GetEmail())
	user := user.User{
		Email:    req.GetEmail(),
		Fullname: req.GetFullname(),
		Password: req.GetPassword(),
	}

	if user.Email == "" || user.Fullname == "" || user.Password == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Fullname, Email, and Password Is Required!",
		)
	}

	u, _ := s.Manager.FindOne(user.Email)
	if u != nil {
		return nil, status.Errorf(
			codes.AlreadyExists,
			fmt.Sprintf("User with email %s already exist", user.Email),
		)
	}

	_ = s.Manager.Register(user)
	token := s.Token.Generate(&user)

	log.Println("Success registering new user with email:", user.Email)
	return &authpb.RegisterResponse{
		Message:     "Register success",
		AccessToken: token,
	}, nil
}

// LoginService Handler
func (s *Server) LoginService(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	user, err := s.Manager.FindOne(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"Username or Email Combination Does Not Match",
			)
		}

		return nil, status.Errorf(
			codes.Internal,
			"Internal Error",
		)
	}

	if ok := pwd.Compare(password, user.Password); !ok {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"Username or Email Combination Does Not Match",
		)
	}

	accessToken := s.Token.Generate(user)

	return &authpb.LoginResponse{
		Message:     "Login Success",
		TokenType:   "Bearer",
		AccessToken: accessToken,
	}, nil
}
