package config

import "os"

type Server struct {
	Port     string
	CertPath string
}

// TODO: Take log as argument for handling empty values
func GetServer() *Server {
	srvCfg := &Server{}
	srvCfg.Port = ":" + os.Getenv("SERVER_PORT")
	srvCfg.CertPath = os.Getenv("SERVER_CERT_PATH")
	return srvCfg
}
