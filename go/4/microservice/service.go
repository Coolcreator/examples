package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	acl         map[string][]string
	defaultHost string

	Admin
	Biz
}

type Admin struct {
	ctx context.Context
	mu  sync.Mutex

	logChannel  chan *Event
	statChannel chan *Event

	logSubscribers  []chan *Event
	statSubscribers []chan *Event
}

func (s *Server) Subscribe(subs *[]chan *Event) chan *Event {
	ch := make(chan *Event)
	s.mu.Lock()
	defer s.mu.Unlock()
	*subs = append(*subs, ch)
	return ch
}

func UpdateStat() *Stat {
	return &Stat{
		ByMethod:   make(map[string]uint64),
		ByConsumer: make(map[string]uint64),
	}
}

func (s *Server) Logging(nothing *Nothing, srv Admin_LoggingServer) error {
	ch := s.Subscribe(&s.logSubscribers)

	for {
		select {
		case event := <-ch:
			srv.Send(event)
		case <-s.ctx.Done():
			return nil
		}
	}
}

func (s *Server) Statistics(interval *StatInterval, srv Admin_StatisticsServer) error {
	ch := s.Subscribe(&s.statSubscribers)

	ticker := time.NewTicker(time.Second * time.Duration(interval.IntervalSeconds))
	stat := UpdateStat()

	for {
		select {
		case event := <-ch:
			stat.ByMethod[event.Method] += 1
			stat.ByConsumer[event.Consumer] += 1
		case <-ticker.C:
			srv.Send(stat)
			stat = UpdateStat()
		case <-s.ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

type Biz struct{}

func (s *Server) Add(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *Server) Check(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *Server) Test(ctx context.Context, in *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func StartMyMicroservice(ctx context.Context, listenAddr string, ACLData string) error {
	s := &Server{
		defaultHost: "127.0.0.1:8080",
		Admin: Admin{
			ctx:         ctx,
			logChannel:  make(chan *Event, 0),
			statChannel: make(chan *Event, 0),
		},
	}

	if err := json.Unmarshal([]byte(ACLData), &s.acl); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(s.UnaryInterceptor),
		grpc.StreamInterceptor(s.StreamInterceptor),
	)

	RegisterBizServer(server, s)
	RegisterAdminServer(server, s)

	go s.Handle()

	go func() {
		err = server.Serve(lis)
		if err != nil {
			log.Fatal("Failed to start server", err)
		}
	}()

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	return nil
}

func (s *Server) Handle() {
	for {
		select {
		case event := <-s.logChannel:
			s.mu.Lock()
			for _, channel := range s.logSubscribers {
				channel <- event
			}
			s.mu.Unlock()
		case stat := <-s.statChannel:
			s.mu.Lock()
			for _, channel := range s.statSubscribers {
				channel <- stat
			}
			s.mu.Unlock()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Server) AccessControl(ctx context.Context, method string) (string, error) {
	consumer, err := GetConsumer(ctx)
	if err != nil {
		return "", err
	}

	methods, ok := s.acl[consumer]
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "Can't get methods")
	}

	if strings.Contains(methods[0], "*") {
		return consumer, nil
	}

	for _, m := range methods {
		if m == method {
			return consumer, nil
		}
	}

	return "", grpc.Errorf(codes.Unauthenticated, "Unauthorized")
}

func (s *Server) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	consumer, err := s.AccessControl(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	event := &Event{
		Consumer: consumer,
		Method:   info.FullMethod,
		Host:     s.defaultHost,
	}

	s.logChannel <- event
	s.statChannel <- event

	return handler(ctx, req)
}

func (s *Server) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	consumer, err := s.AccessControl(ss.Context(), info.FullMethod)
	if err != nil {
		return err
	}

	event := &Event{
		Consumer: consumer,
		Method:   info.FullMethod,
		Host:     s.defaultHost,
	}

	s.logChannel <- event
	s.statChannel <- event

	return handler(srv, ss)
}

func GetConsumer(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "Can't get metadata")
	}

	cons, ok := md["consumer"]
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "Can't get consumer")
	}

	return cons[0], nil
}
