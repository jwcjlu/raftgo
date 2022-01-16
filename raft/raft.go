package raft

import (
	"context"
	"fmt"
	"github.com/raftgo/api"
	"github.com/raftgo/config"
	"github.com/raftgo/store"
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
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	FetchEntries(context.Context, *api.FetchEntryRequest) (*api.FetchEntryResponse, error)
	SendEntry(context.Context, *api.Entry) (*api.Response, error)
	CommitEntry(context.Context, *api.CommitEntryReq) (*api.Response, error)
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
	db          store.Store
	isLeader    bool
	lastIndex   int64
	commitIndex int64
}

func NewRaft(conf *config.Config, db store.Store) Raft {
	raft := raft{stable: NewSet(), heartbeatCh: make(chan *api.HeartbeatRequest, 1)}
	raft.Init(conf)
	raft.db = db
	return &raft
}

func (r *raft) Init(conf *config.Config) {
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
				r.isLeader = false
				r.leaderId = req.NodeId
			}
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 1)
			r.heartbeat(&api.HeartbeatRequest{Term: r.Term, NodeId: r.Id})
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
				r.isLeader = false
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
		r.isLeader = true
	}
}

func (r *raft) heartbeat(request *api.HeartbeatRequest) {

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
			api.NewApiServiceClient(conn).Heartbeat(ctx, request)

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
func (r *raft) Put(key, value []byte) error {
	if !r.isLeader {
		return fmt.Errorf("not leader")
	}
	r.lastIndex++
	entry := store.Entry{
		Key:   key,
		Value: value,
		Meta:  store.Meta{Term: r.Term, LastIndex: r.lastIndex},
	}
	r.db.Put(entry)
	r.sendEntry(entry)
	return nil
}

func (r *raft) sendEntry(entry store.Entry) {
	set := NewSet()
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
			_, err = api.NewApiServiceClient(conn).SendEntry(ctx, &api.Entry{
				Term:   entry.Meta.Term,
				Key:    entry.Key,
				Value:  entry.Value,
				LastId: entry.Meta.LastIndex,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			set.add(node)

		}(node)
	}
	wg.Wait()
	logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Infof("send msg finish:%v", entry)

	if set.len() > 0 {
		r.db.Commit(entry.Meta.Term, entry.Meta.LastIndex)
		r.commitIndex = entry.Meta.LastIndex
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Infof("commit start:%v", entry)
		wg.Add(set.len())
		for _, node := range set.data {
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
				_, err = api.NewApiServiceClient(conn).CommitEntry(ctx, &api.CommitEntryReq{
					Term:    entry.Meta.Term,
					EntryId: entry.Meta.LastIndex,
				})
				if err != nil {
					logrus.Error(err)
					return
				}

				set.add(node)

			}(node)
		}
	}
	wg.Wait()
	logrus.WithField("node", fmt.Sprintf("%s:%d", r.Ip, r.Port)).Infof("commit end:%v", entry)
}

func (r *raft) Get(key []byte) ([]byte, error) {
	return r.db.Get(key)
}

func (r *raft) Leader(ctx context.Context, req *api.LeaderRequest) (*api.Response, error) {
	return nil, nil
}
func (r *raft) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (*api.Response, error) {
	r.heartbeatCh <- req
	return &api.Response{Code: 100000}, nil
}

func (r *raft) FetchEntries(ctx context.Context, req *api.FetchEntryRequest) (*api.FetchEntryResponse, error) {
	if !r.isLeader {
		return &api.FetchEntryResponse{}, fmt.Errorf("not leader")
	}
	entries, err := r.db.FetchEntries(req.Term, req.LastId)
	if err != nil {
		return &api.FetchEntryResponse{}, err
	}
	var data []*api.Entry
	for _, entry := range entries {
		data = append(data, &api.Entry{
			Term:   entry.Meta.Term,
			Key:    entry.Key,
			Value:  entry.Value,
			LastId: entry.Meta.LastIndex,
		})
	}

	return &api.FetchEntryResponse{Entries: data}, nil
}
func (r *raft) SendEntry(ctx context.Context, entry *api.Entry) (*api.Response, error) {
	if entry.Term < r.Term {
		return nil, fmt.Errorf("term is invalid")
	}
	r.db.Put(store.Entry{
		Key:   entry.Key,
		Value: entry.Value,
		Meta:  store.Meta{Term: entry.Term, LastIndex: entry.LastId},
	})
	r.lastIndex = entry.LastId
	return &api.Response{Code: 100000}, nil
}
func (r *raft) CommitEntry(ctx context.Context, req *api.CommitEntryReq) (*api.Response, error) {
	if req.Term < r.Term {
		return nil, fmt.Errorf("term is invalid")
	}
	if req.EntryId > r.lastIndex {
		return nil, fmt.Errorf("catch up")
	}
	r.db.Commit(req.Term, req.EntryId)
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
