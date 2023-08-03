package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	eventpb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *HandlerGRPC) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
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
	return &eventpb.CreateEventResponse{
		Id: id,
	}, nil
}

func (h *HandlerGRPC) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	var date time.Time
	if req.Event.GetDate() != nil {
		date = req.Event.GetDate().AsTime()
	}

	id := req.Event.GetId()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	event := models.Event{
		ID:                   parsedID.String(),
		Title:                req.GetEvent().GetTitle(),
		Date:                 date,
		Duration:             req.GetEvent().GetDuration().AsDuration(),
		Description:          req.GetEvent().GetDescription(),
		UserID:               int(req.GetEvent().GetUserId()),
		NotificationInterval: req.GetEvent().GetNotificationInterval().AsDuration(),
	}

	updatedEvent, err := h.service.UpdateEvent(ctx, event.ID, event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUpdatedEvent := toPBEvent(updatedEvent)
	resultEvent := eventpb.UpdateEventResponse{
		Event: &pbUpdatedEvent,
	}

	return &resultEvent, nil
}

func (h *HandlerGRPC) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*emptypb.Empty, error) {
	parsedID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.service.DeleteEvent(ctx, parsedID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *HandlerGRPC) ListEventsByDay(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) { //nolint:lll
	events, err := h.service.GetAllByDayEvents(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*eventpb.Event, 0, len(events))

	for _, event := range events {
		pbEvent := toPBEvent(event)
		result = append(result, &pbEvent)
	}

	return &eventpb.ListEventsResponse{
		Events: result,
	}, nil
}

func (h *HandlerGRPC) ListEventsByWeek(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) { //nolint:lll
	parsedDate, err := time.Parse(time.RFC3339, req.Date.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := h.service.GetAllByWeekEvents(ctx, parsedDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*eventpb.Event, 0, len(events))

	for _, event := range events {
		pbEvent := toPBEvent(event)
		result = append(result, &pbEvent)
	}

	return &eventpb.ListEventsResponse{
		Events: result,
	}, nil
}

func (h *HandlerGRPC) ListEventsByMonth(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) { //nolint:lll
	parsedDate, err := time.Parse(time.RFC3339, req.Date.String())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := h.service.GetAllByMonthEvents(ctx, parsedDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*eventpb.Event, 0, len(events))

	for _, event := range events {
		pbEvent := toPBEvent(event)
		result = append(result, &pbEvent)
	}

	return &eventpb.ListEventsResponse{
		Events: result,
	}, nil
}

func toPBEvent(event models.Event) eventpb.Event {
	return eventpb.Event{
		Id:                   event.ID,
		Title:                event.Title,
		Date:                 timestamppb.New(event.Date),
		Duration:             durationpb.New(event.Duration),
		Description:          event.Description,
		UserId:               int64(event.UserID),
		NotificationInterval: durationpb.New(event.NotificationInterval),
	}
}
