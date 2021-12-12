package raft

import (
	"context"
	"github.com/raftgo/api"
)

var SuccessFul = &service.Response{Code: 10000}

type Raft interface {
	LifeCycle
	VoteFor(ctx context.Context, request *service.VoteRequest) (*service.VoteResponse, error)
}
type raft struct {
	Id    string
	Ip    string
	Port  int
	Term  int32
	nodes []string
	State StateEnum
}

func NewRaft(conf *Config) Raft {
	raft := raft{}
	raft.Init(conf)
	return &raft
}

func (r *raft) Init(conf *Config) {
	r.nodes = conf.Cluster.Nodes
	r.Id = conf.Node.Id
	r.Ip = conf.Node.Ip
	r.Port = conf.Node.Port

}
func (r *raft) Start() {
	go r.Run()
}

func (r *raft) Stop() {

}
func (r *raft) Run() {
	for {
		switch r.State {
		case Unknown:
			r.State = Candidate
		case Candidate:
			r.runCandidate()
		case Leader:
			r.runLeader()
		case Follower:
			r.runFollower()
		}

	}

}

func (r *raft) runLeader() {

}
func (r *raft) runCandidate() {

}

func (r *raft) runFollower() {

}

func (r *raft) VoteFor(ctx context.Context, request *service.VoteRequest) (*service.VoteResponse, error) {
	if request.Term > r.Term {
		return &service.VoteResponse{Rsp: SuccessFul, Result: true}, nil
	}
	return &service.VoteResponse{Rsp: SuccessFul, Result: false}, nil
}
