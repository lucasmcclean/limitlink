package config

import (
	"os"
	"strings"

	"github.com/lucasmcclean/url-shortener/logger"
)

type Environment int

const (
	DEV Environment = iota
	TEST
	PROD
)

type App struct {
	Env Environment
}

func GetApp(log logger.Logger) *App {
	appCfg := &App{}

	env := os.Getenv("ENVIRONMENT")
	env = strings.ToLower(env)

	switch env {
	case "dev":
		appCfg.Env = DEV
	case "test":
		appCfg.Env = TEST
	case "prod":
		appCfg.Env = PROD
	default:
		log.Warn("an invalid environment was specified; defaulting to PROD", "env", env)
		appCfg.Env = PROD
	}

	return appCfg
}

func (a *App) IsDev() bool  { return a.Env == DEV }
func (a *App) IsTest() bool { return a.Env == TEST }
func (a *App) IsProd() bool { return a.Env == PROD }
