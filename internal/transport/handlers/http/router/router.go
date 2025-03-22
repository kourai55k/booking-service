package router

import "net/http"

type UserHandler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetUserByLogin(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	mux         *http.ServeMux
	userHandler UserHandler
}

func NewRouter(userHandler UserHandler) *Router {
	r := &Router{
		mux:         http.NewServeMux(),
		userHandler: userHandler,
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

	return r.mux
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
