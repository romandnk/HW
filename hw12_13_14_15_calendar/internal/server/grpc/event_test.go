package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	mock_logger "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger/mock"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	mock_service "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func startGRPCServer() (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)

	srv := grpc.NewServer()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	return srv, listener
}

func getDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestHandlerGRPCCreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, lis := startGRPCServer()
	defer srv.Stop()
	defer lis.Close()

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)
	handler := HandlerGRPC{
		service: services,
		logger:  logger,
	}

	event_pb.RegisterEventServiceServer(srv, &handler)

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := event_pb.NewEventServiceClient(conn)

	id := uuid.New().String()
	event := models.Event{
		Title:                "test",
		Date:                 time.Now().UTC(),
		Duration:             time.Second,
		Description:          "test",
		UserID:               1,
		NotificationInterval: time.Second,
	}

	pbEvent := &event_pb.CreateEventRequest{
		Title:                event.Title,
		Date:                 timestamppb.New(event.Date),
		Duration:             durationpb.New(event.Duration),
		Description:          event.Description,
		UserId:               int64(event.UserID),
		NotificationInterval: durationpb.New(event.NotificationInterval),
	}

	services.EXPECT().CreateEvent(gomock.Any(), event).Return(id, nil)

	res, err := client.CreateEvent(ctx, pbEvent)
	require.NoError(t, err)
	require.Equal(t, id, res.GetId())
}

func TestHandlerGRPCCreateEventError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv, lis := startGRPCServer()
	defer srv.Stop()
	defer lis.Close()

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)
	handler := HandlerGRPC{
		service: services,
		logger:  logger,
	}

	event_pb.RegisterEventServiceServer(srv, &handler)

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getDialer(lis)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := event_pb.NewEventServiceClient(conn)

	id := ""
	event := models.Event{
		Title:                "",
		Date:                 time.Now().UTC(),
		Duration:             time.Second,
		Description:          "test",
		UserID:               1,
		NotificationInterval: time.Second,
	}

	pbEvent := &event_pb.CreateEventRequest{
		Title:                event.Title,
		Date:                 timestamppb.New(event.Date),
		Duration:             durationpb.New(event.Duration),
		Description:          event.Description,
		UserId:               int64(event.UserID),
		NotificationInterval: durationpb.New(event.NotificationInterval),
	}

	services.EXPECT().CreateEvent(gomock.Any(), event).Return("", errors.New("title cannot be empty"))

	res, err := client.CreateEvent(ctx, pbEvent)
	expectedErr := "rpc error: code = Internal desc = title cannot be empty"
	require.Equal(t, expectedErr, err.Error())
	require.Equal(t, id, res.GetId())
}
