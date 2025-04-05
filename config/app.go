package config

import (
	"os"
	"strings"

	"github.com/lucasmcclean/url-shortener/logger"
)

// Environment represents the application environment.
type Environment int

const (
	DEV Environment = iota
	TEST
	PROD
)

// app holds global application configuration.
type app struct {
	env Environment
}

// App is the singleton instance of application config.
var App *app

// InitApp initializes the global App configuration.
// It must be called once at startup.
func InitApp(log logger.Logger) {
  if App != nil {
    log.Fatal("global App has already been initialized")
  }

	App = &app{}

	env := os.Getenv("ENVIRONMENT")
	env = strings.ToLower(env)

	switch env {
	case "dev":
		App.env = DEV
	case "test":
		App.env = TEST
	case "prod":
		App.env = PROD
	default:
		log.Warn("an invalid environment was specified; defaulting to PROD", "env", env)
		App.env = PROD
	}
}

// Env returns the current environment.
func (a *app) Env() Environment { return a.env }

// IsDev reports if the environment is development.
func (a *app) IsDev() bool  { return a.env == DEV }

// IsTest reports if the environment is testing.
func (a *app) IsTest() bool { return a.env == TEST }

// IsProd reports if the environment is production.
func (a *app) IsProd() bool { return a.env == PROD }

// String returns the environment as a string.
func (e Environment) String() string {
    switch e {
    case DEV:
        return "dev"
    case TEST:
        return "test"
    case PROD:
        return "prod"
    default:
        return "unknown"
    }
}
