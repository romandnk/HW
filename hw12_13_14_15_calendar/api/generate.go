package api

//go:generate protoc --go_out=../internal/server/grpc/pb/event --go-grpc_out=../internal/server/grpc/pb/event event/EventService.proto
