package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"

	"calendar.com/middleware"
	"calendar.com/pkg/handler"
	pg "calendar.com/proto"
)

type HTTPHandlers struct {
	*handler.Controller
}

type gRPCHandlers struct {
	pg.UnimplementedCalendarServer
}

func RunServer(ctx context.Context, r *mux.Router) error {
	server := &http.Server{Addr: ":8000", Handler: r}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}

func (h *HTTPHandlers) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	s := r.PathPrefix("/api").Subrouter()
	e := s.PathPrefix("/events").Subrouter()
	e.HandleFunc("/{id}", h.Update).Methods(http.MethodPut)
	e.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
	e.Use(middleware.Authorization)
	e.HandleFunc("", h.Create).Methods(http.MethodPost)

	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, _ := route.GetMethods()
		fmt.Println("-----------------------------")
		str := fmt.Sprintf("%s | %s\n", strings.Join(m, ","), t)
		fmt.Println(str)

		return nil
	})

	return r
}

func NewRouter(ctx context.Context) error {
	fmt.Println("grpc server")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		lis.Close()
	}()

	s := grpc.NewServer()
	pg.RegisterCalendarServer(s, &gRPCHandlers{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v\n", lis.Addr())
	}

	return nil
}

func (s *gRPCHandlers) Login(_ context.Context, _ *pg.Credentials) (*pg.Token, error) {
	return &pg.Token{Token: "some token", ExpiresAt: 1234567890}, nil
}

func (h *HTTPHandlers) NewHandler(c handler.Controller) {
	h.Controller = &c
}
