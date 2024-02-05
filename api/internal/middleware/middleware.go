package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/api/internal/log"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

const (
	AUTH_TOKEN = "auth_token"

	ERR_NOT_AUTHORIZED = "not authorized"
	ERR_NOT_ALLOWED    = "not allowed"
)

func CheckIfAuthorized(c *fiber.Ctx) error {
	token := c.Cookies(AUTH_TOKEN)
	l := log.GetLogger()

	if token == "" {
		c.Status(http.StatusMethodNotAllowed)
		b, err := json.Marshal(map[string]string{"error": ERR_NOT_AUTHORIZED})
		if err != nil {
			l.Error("middleware: failed to marshal error message", slog.String("error", err.Error()))
			return nil
		}

		if err = c.Send(b); err != nil {
			l.Error("middleware: failed to send response", slog.String("error", err.Error()))
		}
		return nil
	}

	c.Context().SetUserValue(AUTH_TOKEN, token)
	fmt.Println(c.UserContext().Value(AUTH_TOKEN))

	if err := c.Next(); err != nil {
		l.Error("failed to execute next method", slog.String("error", err.Error()))
	}

	return nil
}
