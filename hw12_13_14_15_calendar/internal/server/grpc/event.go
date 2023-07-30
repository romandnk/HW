package grpc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var (
	ErrNeedID        = errors.New("there must be id")
	ErrEmptyID       = errors.New("id cannot be empty")
	ErrNeedDate      = errors.New("there must be date")
	ErrEmptyDate     = errors.New("date cannot be empty")
	ErrEmptyPeriod   = errors.New("period cannot be empty")
	ErrInvalidPeriod = errors.New("period can be \"day\", \"week\", \"month\"")
)

func (h *HandlerGRPC) CreateEvent(ctx context.Context, req *event_pb.Event) (*event_pb.CreateEventResponse, error) {
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

func (h *HandlerGRPC) UpdateEvent(ctx context.Context, req *event_pb.UpdateEventRequest) (*event_pb.Event, error) {
	panic("")
}

func (h *HandlerGRPC) DeleteEvent(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, ErrNeedID.Error())
	}

	ids := md.Get("id")
	if len(ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrEmptyID.Error())
	}

	id := ids[0]
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.service.DeleteEvent(ctx, parsedID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *HandlerGRPC) ListEvents(ctx context.Context, _ *emptypb.Empty) (*event_pb.ListEventsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, ErrNeedDate.Error())
	}

	periods := md.Get("period")
	if len(periods) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrEmptyPeriod.Error())
	}

	dates := md.Get("date")
	if len(dates) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrEmptyDate.Error())
	}

	date := dates[0]
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var events []models.Event
	switch periods[0] {
	case "day":
		events, err = h.service.GetAllByDayEvents(ctx, parsedDate)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "week":
		events, err = h.service.GetAllByWeekEvents(ctx, parsedDate)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "month":
		events, err = h.service.GetAllByMonthEvents(ctx, parsedDate)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	default:
		return nil, status.Error(codes.InvalidArgument, ErrInvalidPeriod.Error())
	}

	result := make([]*event_pb.Event, 0, len(events))
	for _, event := range events {
		pbEvent := &event_pb.Event{
			Id:                   &event.ID,
			Title:                event.Title,
			Date:                 timestamppb.New(event.Date),
			Duration:             durationpb.New(event.Duration),
			Description:          event.Description,
			UserId:               int64(event.UserID),
			NotificationInterval: durationpb.New(event.NotificationInterval),
		}
		result = append(result, pbEvent)
	}

	return &event_pb.ListEventsResponse{
		Events: result,
	}, nil
}
