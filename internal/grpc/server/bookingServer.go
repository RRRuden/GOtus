package grpc

import (
	"context"
	"gotus/internal/grpc/api/booking_api"
	"gotus/internal/service/booking"
	"log"
	"net"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

type bookingServer struct {
	booking_api.UnimplementedBookingServiceServer
	service booking.Service
}

var (
	grpcServer *grpc.Server
	listener   net.Listener
)

func NewBookingServer(s booking.Service) *bookingServer {
	return &bookingServer{service: s}
}

func (s *bookingServer) CreateBooking(ctx context.Context, req *booking_api.CreateBookingRequest) (*booking_api.CreateBookingResponse, error) {
	reservation, err := s.service.CreateBooking(int(req.UserId), req.Isbn)
	if err != nil {
		return nil, err
	}
	return &booking_api.CreateBookingResponse{BookingId: int32(reservation.GetID())}, nil
}

func (s *bookingServer) ExtendBooking(ctx context.Context, req *booking_api.ExtendBookingRequest) (*emptypb.Empty, error) {
	_, err := s.service.ExtendBooking(int(req.BookingId), int(req.ExtensionDays))
	return &emptypb.Empty{}, err
}

func (s *bookingServer) CancelBooking(ctx context.Context, req *booking_api.CancelBookingRequest) (*emptypb.Empty, error) {
	_, err := s.service.CancelBooking(int(req.BookingId))
	return &emptypb.Empty{}, err
}

func (s *bookingServer) EndBooking(ctx context.Context, req *booking_api.EndBookingRequest) (*emptypb.Empty, error) {
	_, err := s.service.EndBooking(int(req.BookingId))
	return &emptypb.Empty{}, err
}

func RunGRPCServer(service booking.Service, listenAddr string) error {
	var err error
	listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	grpcServer = grpc.NewServer()
	booking_api.RegisterBookingServiceServer(grpcServer, NewBookingServer(service))

	log.Printf("gRPC сервер запущен на %s", listenAddr)
	return grpcServer.Serve(listener)
}

func StopGRPCServer() {
	if grpcServer != nil {
		log.Println("Остановка gRPC сервера...")
		grpcServer.GracefulStop()
	}
}
