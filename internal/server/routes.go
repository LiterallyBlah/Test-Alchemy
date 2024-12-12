package server

import (
	"net/http"

	"TestAlchemy/cmd/web"
	"TestAlchemy/internal/database"
	"TestAlchemy/internal/handlers"
	"TestAlchemy/internal/middleware"
	"TestAlchemy/internal/services"
	"TestAlchemy/internal/session"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize services and handlers
	db := database.New()
	sessionStore, err := session.NewStore()
	if err != nil {
		panic(err)
	}
	userService := services.NewUserService(db, sessionStore)
	userHandler := handlers.NewUserHandler(userService)

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	// Public routes
	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))
	e.POST("/api/register", echo.WrapHandler(http.HandlerFunc(userHandler.Register)))
	e.POST("/api/login", echo.WrapHandler(http.HandlerFunc(userHandler.Login)))
	e.GET("/health", s.healthHandler)

	// Protected routes - require authentication
	protected := e.Group("")
	protected.Use(middleware.RequireAuth(sessionStore))
	protected.GET("/", s.HelloWorldHandler)
	// Add other protected routes here

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
