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
	conf *Config
}

func NewServer(raft Raft, conf *Config) *Server {
	return &Server{raft: raft, conf: conf}
}

func (server *Server) Vote(ctx context.Context, req *api.VoteRequest) (*api.VoteResponse, error) {

	return server.raft.VoteFor(ctx, req)
}
func (server *Server) Leader(ctx context.Context, req *api.LeaderRequest) (*api.Response, error) {
	return server.raft.Leader(ctx, req)
}
func (server *Server) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (*api.Response, error) {
	return server.raft.Heartbeat(ctx, req)
}

type ServiceManager struct {
	Server *Server
}

func NewServiceManager(server *Server) *ServiceManager {
	return &ServiceManager{Server: server}
}
func (manager ServiceManager) RegisterService(server *grpc.Server) {
	api.RegisterApiServiceServer(server, manager.Server)
}

func (manager *ServiceManager) Start() (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", manager.Server.conf.Node.Port))
	if err != nil {
		return nil, err
	}
	manager.Server.raft.Start()
	return listener, err
}
