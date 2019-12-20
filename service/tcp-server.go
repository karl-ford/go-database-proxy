package service

type TCPServer struct {
	*ServerBase
}

func CreateTCPServer(base *ServerBase) *TCPServer{
	server := &TCPServer{
		base,
	}
	return server
}

func (server *TCPServer) Run() error {
	return nil
}