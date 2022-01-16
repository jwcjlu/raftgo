package raft

import (
	"context"
	"fmt"
	"github.com/raftgo/config"
	"net"

	"google.golang.org/grpc"

	"github.com/raftgo/api"
)

type Server struct {
	raft Raft
	conf *config.Config
}

func NewServer(raft Raft, conf *config.Config) *Server {
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
func (server *Server) FetchEntries(ctx context.Context, req *api.FetchEntryRequest) (*api.FetchEntryResponse, error) {
	return server.raft.FetchEntries(ctx, req)
}
func (server *Server) SendEntry(ctx context.Context, req *api.Entry) (*api.Response, error) {
	return server.raft.SendEntry(ctx, req)
}
func (server *Server) CommitEntry(ctx context.Context, req *api.CommitEntryReq) (*api.Response, error) {
	return server.raft.CommitEntry(ctx, req)
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
