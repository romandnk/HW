package grpc

import (
	"context"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *HandlerGRPC) CreateEvent(ctx context.Context, req *event_pb.CreateEventRequest) (*event_pb.CreateEventResponse, error) {
	event := models.Event{
		Title:                req.GetTitle(),
		Date:                 req.GetDate().AsTime(),
		Duration:             req.GetDuration().AsDuration(),
		Description:          req.GetDescription(),
		UserID:               int(req.GetUserId()),
		NotificationInterval: req.GetNotificationInterval().AsDuration(),
	}

	id, err := h.service.CreateEvent(ctx, event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &event_pb.CreateEventResponse{
		Id: id,
	}, nil
}

func (h *HandlerGRPC) UpdateEvent(context.Context, *event_pb.UpdateEventRequest) (*event_pb.Event, error) {
	panic("")
}

func (h *HandlerGRPC) DeleteEvent(context.Context, *event_pb.DeleteEventRequest) (*emptypb.Empty, error) {
	panic("")
}

func (h *HandlerGRPC) ListEventsByDay(context.Context, *event_pb.ListEventsRequest) (*event_pb.ListEventsResponse, error) {
	panic("")
}

func (h *HandlerGRPC) ListEventsByWeek(context.Context, *event_pb.ListEventsRequest) (*event_pb.ListEventsResponse, error) {
	panic("")
}

func (h *HandlerGRPC) ListEventsByMonth(context.Context, *event_pb.ListEventsRequest) (*event_pb.ListEventsResponse, error) {
	panic("")
}
