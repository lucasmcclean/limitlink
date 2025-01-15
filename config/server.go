package config

import "os"

type Server struct {
	Port string
}

// TODO: Take log as argument for handling empty values
func GetServer() *Server {
	srvCfg := &Server{}
	srvCfg.Port = os.Getenv("SERVER_PORT")
	return srvCfg
}
