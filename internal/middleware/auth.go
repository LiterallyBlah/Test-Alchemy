package middleware

import (
	"TestAlchemy/internal/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequireAuth(sessionStore *session.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			sess, err := sessionStore.GetSession(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid session"})
			}
			if sess == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "session expired"})
			}

			// Store user ID in context for later use
			c.Set("user_id", sess.UserID)
			return next(c)
		}
	}
}
