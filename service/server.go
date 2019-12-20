package service

import (
	"proxy/database"
	"strings"
)

type ServerBase struct {
	port *uint16
	pass *string
	dialects map[string]database.Dialect
}

type Server interface {
	Run() error
	SetDBConnection(name string, dialect database.Dialect) error
}

func CreateServer(port *uint, protocol, pass *string) Server {
	var (
		sport = uint16(*port)
		server = &ServerBase{
			port: &sport,
			pass: pass,
		}
	)

	*protocol = strings.ToLower(*protocol)
	switch *protocol {
	case "http":
		return CreateHTTPServer(server)
	case "tcp":
		return CreateTCPServer(server)
	case "udp":
		return CreateUDPServer(server)
	}

	return nil
}

func (server *ServerBase) SetDBConnection(name string, dialect database.Dialect) error{
	server.dialects[name] = dialect
	return nil
}