package raft

import (
	"context"
	"fmt"
	"github.com/raftgo/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	"time"
)

var SuccessFul = &api.Response{Code: 10000}

type Raft interface {
	LifeCycle
	VoteFor(ctx context.Context, request *api.VoteRequest) (*api.VoteResponse, error)
	Leader(ctx context.Context, in *api.LeaderRequest) (*api.Response, error)
	Heartbeat(ctx context.Context, in *api.HeartbeatRequest) (*api.Response, error)
}
type raft struct {
	Id          string
	Ip          string
	Port        int
	Term        int32
	nodes       []string
	State       StateEnum
	timeout     int
	stable      set
	heartbeatCh chan *api.HeartbeatRequest
	leaderId    string
}

func NewRaft(conf *Config) Raft {
	raft := raft{stable: NewSet(), heartbeatCh: make(chan *api.HeartbeatRequest, 1)}
	raft.Init(conf)
	return &raft
}

func (r *raft) Init(conf *Config) {
	r.nodes = conf.Cluster.Nodes
	r.Id = conf.Node.Id
	r.Ip = conf.Node.Ip
	r.Port = conf.Node.Port
	r.timeout = conf.Node.Timeout

}
func (r *raft) Start() {
	r.State = Follower
	go r.Run()
}

func (r *raft) Stop() {

}
func (r *raft) Run() {
	for {
		switch r.State {
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
	timeoutCh := randomTimeout(time.Second * 1)
	for r.State == Leader {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Info("runLeader")
		select {
		case req := <-r.heartbeatCh:
			if req.Term > r.Term {
				r.State = Follower
				r.leaderId = req.NodeId
			}
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 1)
			r.heartbeat(&api.LeaderRequest{Term: r.Term, NodeId: r.Id})
		}

	}

}
func (r *raft) runCandidate() {
	timeoutCh := randomTimeout(time.Second * 2)
	lastTime := time.Now()
	first := true
	for r.State == Candidate {
		fmt.Println("runCandidate")
		select {
		case req := <-r.heartbeatCh:
			if req.Term > r.Term {
				r.State = Follower
				r.leaderId = req.NodeId
			}
		case <-timeoutCh:
			if time.Now().Sub(lastTime).Seconds() < 2*time.Second.Seconds() && !first {
				first = false
				continue
			}
			timeoutCh = randomTimeout(time.Second * 2)
			r.Term++
			r.sendVoteRequestToAll(&api.VoteRequest{Term: r.Term, NodeId: r.Id})
		}
	}

}

func (r *raft) sendVoteRequestToAll(request *api.VoteRequest) {
	r.stable = NewSet()
	wg := sync.WaitGroup{}
	wg.Add(len(r.nodes))
	for _, node := range r.nodes {
		if node == fmt.Sprintf("%s:%d", r.Ip, r.Port) {
			continue
			wg.Done()
		}
		go func(node string) {
			conn, err := grpc.Dial(node, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer func() {
				if conn != nil {
					conn.Close()
				}
				wg.Done()
			}()
			if err != nil {
				logrus.Error(err)
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			rsp, err := api.NewApiServiceClient(conn).Vote(ctx, request)
			if err != nil {
				logrus.Error(err)
				return
			}
			if rsp.Result {
				r.stable.add(node)
			}

		}(node)
	}
	wg.Wait()
	logrus.Info("wait!")
	if r.stable.len() > 0 {
		r.State = Leader
	}
}

func (r *raft) heartbeat(request *api.LeaderRequest) {

	for _, node := range r.nodes {
		if node == fmt.Sprintf("%s:%d", r.Ip, r.Port) {
			continue
		}
		go func(node string) {
			conn, err := grpc.Dial(node, grpc.WithInsecure())
			defer func() {
				conn.Close()
			}()
			if err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			api.NewApiServiceClient(conn).Leader(ctx, request)

		}(node)
	}
}
func (r *raft) runFollower() {
	timeoutCh := randomTimeout(time.Second * 2)
	lastTime := time.Now()
	for r.State == Follower {
		fmt.Println("runFollower")
		select {
		case req := <-r.heartbeatCh:
			if req.Term > r.Term {
				r.leaderId = req.NodeId
			}
			lastTime = time.Now()
			logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Infof("leader id=%s", r.leaderId)
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 2)
			if time.Now().Sub(lastTime).Seconds() < 2*time.Second.Seconds() {
				continue
			}
			r.State = Candidate
		}
	}
}

func (r *raft) VoteFor(ctx context.Context, request *api.VoteRequest) (*api.VoteResponse, error) {
	logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Infof("requestTerm:%d->currentTerm:%d", request.Term, r.Term)
	if request.Term > r.Term {
		return &api.VoteResponse{Rsp: SuccessFul, Result: true}, nil
	}
	return &api.VoteResponse{Rsp: SuccessFul, Result: false}, nil
}

func (r *raft) Leader(ctx context.Context, req *api.LeaderRequest) (*api.Response, error) {
	return nil, nil
}
func (r *raft) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (*api.Response, error) {
	r.heartbeatCh <- req
	return &api.Response{Code: 100000}, nil
}
func randomTimeout(minVal time.Duration) <-chan time.Time {
	if minVal == 0 {
		return nil
	}
	extra := (time.Duration(rand.Int63()) % minVal)
	return time.After(minVal + extra)
}

type set struct {
	data map[string]string
}

func NewSet() set {
	return set{data: make(map[string]string, 0)}
}
func (s set) add(element string) {
	s.data[element] = element
}
func (s set) remote(element string) {
	delete(s.data, element)
}

func (s set) len() int {
	return len(s.data)
}
