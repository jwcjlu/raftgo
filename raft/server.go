package raft

import (
	"context"
	"fmt"
	"github.com/jwcjlu/raftgo/config"
	"google.golang.org/grpc"
	"net"

	"github.com/jwcjlu/raftgo/api"
)

func NewRaftService(raft *Raft, conf *config.Config) api.RaftServiceServer {
	return &RaftService{raft: raft, conf: conf}
}

type RaftService struct {
	conf *config.Config
	raft *Raft
}

func (r *RaftService) Vote(ctx context.Context, req *api.VoteRequest) (*api.VoteResponse, error) {
	r.raft.voteChan <- req
	rsp := <-r.raft.voteResponseChan
	return rsp, nil
}

// 定义方法
func (r *RaftService) AppendEntries(ctx context.Context, req *api.AppendEntriesRequest) (*api.AppendEntriesResponse, error) {
	r.raft.appendEntryChan <- req
	rsp := <-r.raft.appendEntryResponseChan
	return rsp, nil

}

// 安装快照
func (r *RaftService) InstallSnapshot(ctx context.Context, req *api.InstallSnapshotRequest) (*api.InstallSnapshotResponse, error) {
	return nil, nil
}

type ServiceManager struct {
	Server *RaftService
}

func NewServiceManager(server *RaftService) *ServiceManager {
	return &ServiceManager{Server: server}
}
func (manager ServiceManager) RegisterService(server *grpc.Server) {
	api.RegisterRaftServiceServer(server, manager.Server)
}

func (manager *ServiceManager) Start() (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", manager.Server.conf.Node.Port))
	if err != nil {
		return nil, err
	}
	manager.Server.raft.Start()
	return listener, err
}
