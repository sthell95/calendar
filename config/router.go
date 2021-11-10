package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"calendar.com/pkg/domain/entity"

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

func checkUserRole(next http.Handler) http.Handler {
	fmt.Println("Check user role")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("Check user role inside")

		ctx := context.WithValue(r.Context(), "user", "User")
		ctx2 := context.WithValue(ctx, entity.CtxUserKey2, "User 2")
		ctx3 := context.WithValue(ctx2, entity.CtxUserKey, "User const")
		r = r.WithContext(ctx3)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		fmt.Println("Check authorization")
		return next
	})
	r.Use(checkUserRole)
	r.Use(func(next http.Handler) http.Handler {
		fmt.Println("Check user permissions")
		return next
	})

	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/health_checker", h.HealthHandler).Methods(http.MethodGet)

	s.HandleFunc("/login", h.SignIn).Methods(http.MethodPost)

	s.HandleFunc("/events", h.Create).Methods(http.MethodPost)

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
