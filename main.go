package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/jwt"
	authpb "github.com/mhdiiilham/oauth2-auth-server-implementation/protos"
	"github.com/mhdiiilham/oauth2-auth-server-implementation/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/entity/user"

	"github.com/mhdiiilham/oauth2-auth-server-implementation/mongodb"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, coll, mongoErr := mongodb.NewMongoDBConnection(os.Getenv("MONGO_DB_USER"), os.Getenv("MONGO_DB_PASS"), os.Getenv("MONGO_DB"), os.Getenv("MONGO_DB_COLLECTION"))
	if mongoErr != nil {
		panic(mongoErr)
	}

	tokenService := jwt.NewJWTService(string(os.Getenv("JWT_SECRET")), string(os.Getenv("APP_NAME")))
	userRepo := user.NewMongoDBRepository(coll)
	userManager := user.NewManager(userRepo)

	authServer := service.NewService(os.Getenv("SERVER_NETWORK"), os.Getenv("SERVER_ADDRESS"), userManager, tokenService)

	lis, err := net.Listen(authServer.Network, authServer.Address)
	if err != nil {
		log.Printf("Failed listening to port 50051 due to error: %v", err.Error())
		panic(err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	authpb.RegisterAuthorizationServiceServer(s, authServer)

	go func() {
		log.Println("Starting authorization service")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Fail to serve: Error: %v", err.Error())
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Stopping the server")
	s.Stop()
	log.Println("Closing listener")
	lis.Close()
	log.Println("Disconnection from MongoDB")
	client.Disconnect(context.TODO())
}
