package router

import (
	"net/http"

	"github.com/kourai55k/booking-service/internal/transport/handlers/http/middleware"
)

type UserHandler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetUserByLogin(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ProtectedHello(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	mux         *http.ServeMux
	userHandler UserHandler
	authHandler AuthHandler
}

func NewRouter(userHandler UserHandler, authHandler AuthHandler) *Router {
	r := &Router{
		mux:         http.NewServeMux(),
		userHandler: userHandler,
		authHandler: authHandler,
	}
	r.RegisterRoutes()
	return r
}

func (r *Router) RegisterRoutes() *http.ServeMux {
	// users routes
	r.mux.HandleFunc("GET /user/{id}", r.userHandler.GetUserByID)
	r.mux.HandleFunc("GET /user", r.userHandler.GetUserByLogin)
	r.mux.HandleFunc("GET /users", r.userHandler.GetUsers)
	r.mux.HandleFunc("POST /user", r.userHandler.CreateUser)
	r.mux.HandleFunc("PATCH /user/{id}", r.userHandler.UpdateUser)
	r.mux.HandleFunc("DELETE /user/{id}", r.userHandler.DeleteUser)

	// auth routes
	r.mux.HandleFunc("/register", r.authHandler.Register)
	r.mux.HandleFunc("/login", r.authHandler.Login)

	// test route for testing middleware
	r.mux.Handle("/protected/hello", middleware.AuthMiddleware(http.HandlerFunc(r.userHandler.ProtectedHello)))

	// restrants routes

	// bookings routes

	return r.mux
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
