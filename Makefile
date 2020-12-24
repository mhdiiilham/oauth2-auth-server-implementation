genproto:
	protoc protos/auth.proto --go_out=plugins=grpc:.