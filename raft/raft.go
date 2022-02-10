package raft

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jwcjlu/raftgo/api"
	"github.com/jwcjlu/raftgo/config"
	"github.com/sirupsen/logrus"
)

//--------------------------------------
// Event Loop
//--------------------------------------

//               ________
//            --|Snapshot|                 timeout
//            |  --------                  ______
// recover    |       ^                   |      |
// snapshot / |       |snapshot           |      |
// higher     |       |                   v      |     recv majority votes
// term       |    --------    timeout    -----------                        -----------
//            |-> |Follower| ----------> | Candidate |--------------------> |  Leader   |
//                 --------               -----------                        -----------
//                    ^          higher term/ |                         higher term |
//                    |            new leader |                                     |
//                    |_______________________|____________________________________ |
// The main event loop for the server
type Raft struct {
	id                      string
	leaderId                string
	ip                      string
	port                    int
	term                    int32
	isLeader                bool
	lastApplied             int64
	commitIndex             int64
	custer                  []*Node
	state                   StateEnum
	appendEntryChan         chan *api.AppendEntriesRequest
	appendEntryResponseChan chan *api.AppendEntriesResponse
	appliedEntryChan        chan *api.LogEntry
	appliedEntryResp        chan *api.LogEntry
	LifeCycle
	mu  sync.RWMutex
	log *Log
}

var timeout = 3

func (r *Raft) Start() {
	r.state = Follower
	go r.Run()
}

func (r *Raft) Stop() {

}
func NewRaft(conf *config.Config) *Raft {
	raft := Raft{}
	raft.Init(conf)
	return &raft
}
func (r *Raft) Init(conf *config.Config) {
	r.mu = sync.RWMutex{}
	r.ip = conf.Node.Ip
	r.port = conf.Node.Port
	r.appendEntryResponseChan = make(chan *api.AppendEntriesResponse, 1)
	r.appendEntryChan = make(chan *api.AppendEntriesRequest, 1)
	for _, ipPort := range conf.Cluster.Nodes {
		ipPorts := strings.Split(ipPort, ":")
		port, _ := strconv.Atoi(ipPorts[1])
		node := &Node{
			Ip:   ipPorts[0],
			Port: port}
		node.raft = r
		node.Init(conf.Node.ConnCount)
		r.custer = append(r.custer, node)
	}
	r.log = &Log{}
	err := r.log.Init(conf)
	if err != nil {
		logrus.Fatal("log file open error =", err)
	}

}

func (r *Raft) Run() {

	for {
		switch r.state {
		case Leader:
			r.runLeader()
		case Follower:
			r.runFollower()
		case Candidate:
			r.runCandidate()
		}
	}

}

func (r *Raft) runLeader() {
	timeoutCh := randomTimeout(time.Millisecond * 800)
	for r.state == Leader {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Info("runLeader")
		select {
		case req := <-r.appendEntryChan:
			rsp := &api.AppendEntriesResponse{Term: r.term, Success: false}
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				rsp.Success = true
			}
			r.appendEntryResponseChan <- rsp
		case <-timeoutCh:
			for _, n := range r.custer {
				go func(node *Node) {
					rsp, err := node.AppendEntries(context.Background(), &api.AppendEntriesRequest{
						Term:     r.term,
						LeaderId: r.leaderId,
					})
					if err != nil {
						logrus.Error(err)
						return
					}
					if rsp.Term > r.term {
						r.state = Follower
						r.isLeader = false
					}
				}(n)
			}
			timeoutCh = randomTimeout(time.Millisecond * 800)
		}
	}
}

func (r *Raft) HandlerVote(request *api.VoteRequest) (*api.VoteResponse, error) {
	rsp := api.VoteResponse{Term: r.term, VoteGranted: false}
	if request.Term > r.term {
		rsp.VoteGranted = true
	}
	return &rsp, nil
}

func (r *Raft) runFollower() {
	timeoutCh := randomTimeout(time.Second * 3)
	lastTime := time.Now()
	for r.state == Follower {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Info("runFollower")
		select {
		case req := <-r.appendEntryChan:
			lastTime = time.Now()
			rsp := &api.AppendEntriesResponse{Term: r.term, Success: false}
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
				rsp.Success = true
			}
			r.appendEntryResponseChan <- rsp

		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 3)
			if time.Now().Sub(lastTime).Seconds() < 3*time.Second.Seconds() {
				continue
			}
			r.state = Candidate

		}
	}
}

func (r *Raft) runCandidate() {
	timeoutCh := randomTimeout(time.Second * 3)
	isVote := true
	for r.state == Candidate {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Infof("runCandidate term=%d", r.term)
		select {
		case req := <-r.appendEntryChan:
			timeoutCh = randomTimeout(time.Second * 3)
			rsp := &api.AppendEntriesResponse{Term: r.term, Success: false}
			if req.Term >= r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
				rsp.Success = true
			}
			r.appendEntryResponseChan <- rsp
		case <-timeoutCh:
			if !isVote {
				timeoutCh = randomTimeout(time.Second * 3)
				continue
			}
			isVote = false
			voteGranted := 1
			r.term++
			wg := sync.WaitGroup{}
			wg.Add(len(r.custer))
			for _, n := range r.custer {
				go func(node *Node, term int32) {
					defer wg.Done()
					rsp, err := node.Vote(context.Background(), &api.VoteRequest{
						Term:         term,
						CandidateId:  r.id,
						LastLogIndex: 0,
						LastLogTerm:  0,
					})
					if err != nil {
						logrus.Error(err)
						return
					}
					if rsp.Term > r.term {
						r.term = rsp.Term
					}
					if rsp.VoteGranted {
						voteGranted++
					}
				}(n, r.term)
			}
			wg.Wait()
			if voteGranted >= r.QuorumSize() {
				r.state = Leader
			}
			timeoutCh = randomTimeout(time.Second * 3)
			isVote = true
		}
	}
}

func (r *Raft) QuorumSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return (len(r.custer) / 2) + 1
}

func randomTimeout(minVal time.Duration) <-chan time.Time {
	if minVal == 0 {
		return nil
	}
	extra := time.Duration(rand.Int63()) % minVal
	return time.After(minVal + extra)
}
