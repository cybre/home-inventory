package toast

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cybre/home-inventory/internal/logging"
	"github.com/labstack/echo/v4"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "danger"
)

var levelTitles = map[string]string{
	INFO:    "Info",
	SUCCESS: "Success",
	WARNING: "Warning",
	DANGER:  "Error",
}

type Toast struct {
	Level      string `json:"level"`
	LevelTitle string `json:"levelTitle"`
	Message    string `json:"message"`
}

func New(level string, message string) Toast {
	return Toast{level, levelTitles[level], message}
}

func Info(message string) Toast {
	return New(INFO, message)
}

func Success(c echo.Context, message string) {
	New(SUCCESS, message).SetHXTriggerHeader(c)
}

func Warning(message string) Toast {
	return New(WARNING, message)
}

func Danger(message string) Toast {
	return New(DANGER, message)
}

func (t Toast) Error() string {
	return fmt.Sprintf("%s: %s", t.Level, t.Message)
}

func (t Toast) jsonify() (string, error) {
	eventMap := map[string]Toast{}
	eventMap["show-toast"] = t
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (t Toast) SetHXTriggerHeader(c echo.Context) {
	jsonData, err := t.jsonify()
	if err != nil {
		logging.FromContext(c.Request().Context()).Warn("failed to jsonify toast", slog.Any("error", err))
		c.Response().Header().Set("HX-Trigger", "{\"show-toast\": {\"level\":\"danger\",\"message\":\"There has been an unexpected error\"}}")
		return
	}

	c.Response().Header().Set("HX-Trigger", jsonData)
}
