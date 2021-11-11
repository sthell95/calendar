package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"calendar.com/middleware"

	"calendar.com/pkg/controller"

	"github.com/gorilla/mux"
)

type Handlers struct {
	*controller.Controller
}

func Run(ctx context.Context, r *mux.Router) error {
	server := &http.Server{Addr: ":8000", Handler: r}

	go func() {
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}

func (h *Handlers) NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/events", h.Create).Methods(http.MethodPost)
	s.Use(middleware.Authorization)

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

func (h *Handlers) NewHandler(c controller.Controller) {
	h.Controller = &c
}
