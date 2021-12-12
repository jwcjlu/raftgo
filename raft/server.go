package raft

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/raftgo/api"
)

type Server struct {
	raft Raft
	conf Config
}

func NewServer(raft Raft, conf Config) *Server {
	return &Server{raft: raft}
}

func (server *Server) Vote(ctx context.Context, req *service.VoteRequest) (*service.VoteResponse, error) {

	return server.raft.VoteFor(ctx, req)
}

type ServiceManager struct {
	Server *Server
}

func NewServiceManager(server *Server) ServiceManager {
	return ServiceManager{Server: server}
}
func (manager ServiceManager) RegisterService(server *grpc.Server) {
	service.RegisterApiServiceServer(server, manager.Server)
}

func (manager *ServiceManager) Start() (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", manager.Server.conf.Node.Port))
	if err != nil {
		return nil, err
	}
	manager.Server.raft.Start()
	return listener, err
}
