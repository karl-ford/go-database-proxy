package service

type UDPServer struct {
	*ServerBase
}

func CreateUDPServer(base *ServerBase) *UDPServer{
	server := &UDPServer{
		base,
	}
	return server
}

func (server *UDPServer) Run() error {
	return nil
}