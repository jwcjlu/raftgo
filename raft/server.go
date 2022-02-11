package raft

import (
	"context"
	"fmt"
	"net"

	"github.com/jwcjlu/raftgo/config"
	"github.com/jwcjlu/raftgo/proto"

	"google.golang.org/grpc"
)

func NewService(raft *Raft, conf *config.Config) *Service {
	return &Service{raft: raft, conf: conf}
}

type Service struct {
	conf *config.Config
	raft *Raft
}

func (r *Service) Vote(ctx context.Context, req *proto.VoteRequest) (*proto.VoteResponse, error) {
	return r.raft.HandlerVote(req)
}

// 定义方法
func (r *Service) AppendEntries(ctx context.Context, req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	r.raft.appendEntryChan <- req
	rsp := <-r.raft.appendEntryResponseChan
	return rsp, nil

}

// 定义方法
func (r *Service) Heartbeat(ctx context.Context, req *proto.HeartbeatRequest) (*proto.HeartbeatResponse, error) {
	r.raft.heartbeatChan <- req
	rsp := <-r.raft.heartbeatResponseChan
	return rsp, nil

}

// 安装快照
func (r *Service) InstallSnapshot(ctx context.Context, req *proto.InstallSnapshotRequest) (*proto.InstallSnapshotResponse, error) {
	return nil, nil
}

type ServiceManager struct {
	Server *Service
}

func NewServiceManager(server *Service) *ServiceManager {
	return &ServiceManager{Server: server}
}
func (manager ServiceManager) RegisterService(server *grpc.Server) {
	proto.RegisterRaftServiceServer(server, manager.Server)
}

func (manager *ServiceManager) Start() (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", manager.Server.conf.Node.Port))
	if err != nil {
		return nil, err
	}
	manager.Server.raft.Start()
	return listener, err
}
