package config

import (
	"github.com/lucasmcclean/url-shortener/logger"
)

type Server struct {
	Port     string
	CertPath string
	Env      string
}

func GetServer(log logger.Logger) *Server {
	srvCfg := &Server{}
	var missing []string

	srvCfg.Port, missing = getOrAppendMissing("SERVER_PORT", missing)
	srvCfg.CertPath, missing = getOrAppendMissing("SERVER_CERT_PATH", missing)
	srvCfg.Env, missing = getOrAppendMissing("ENVIRONMENT", missing)

	if len(missing) > 0 {
		log.Fatal("missing server environment variables\n", "missing_env_vars", missing)
	}

	srvCfg.Port = ":" + srvCfg.Port

	return srvCfg
}
