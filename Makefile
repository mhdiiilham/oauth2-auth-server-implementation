genproto:
	protoc protos/auth.proto --go_out=plugins=grpc:.
test:
	go test ./...
securitycheck:
	go get github.com/securego/gosec/v2/cmd/gosec
	gosec -exclude=G104 -tests ./...