package raft

import (
	"context"
	"fmt"
	"github.com/jwcjlu/raftgo/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jwcjlu/raftgo/api"
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
	term                    int32
	isLeader                bool
	lastApplied             int64
	commitIndex             int64
	custer                  []*Node
	state                   StateEnum
	voteChan                chan *api.VoteRequest
	appendEntryChan         chan *api.AppendEntriesRequest
	voteResponseChan        chan *api.VoteResponse
	appendEntryResponseChan chan *api.AppendEntriesResponse
	LifeCycle
}

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
	for _, ipPort := range conf.Cluster.Nodes {
		ipPorts := strings.Split(ipPort, ":")
		port, _ := strconv.Atoi(ipPorts[1])
		r.custer = append(r.custer, &Node{
			Ip:    ipPorts[0],
			Port:  port,
			State: 0,
			pool: Pool{factory: func() (io.Closer, error) {
				conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ipPorts[0], port),
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					return nil, err
				}
				return conn, nil
			}},
		})
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
	timeoutCh := randomTimeout(time.Microsecond * 100)
	for r.state == Leader {
		select {
		case req := <-r.appendEntryChan:
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
			}
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 2)
		}
	}
}

func (r *Raft) runFollower() {
	timeoutCh := randomTimeout(time.Second * 2)
	for r.state == Follower {
		select {
		case req := <-r.appendEntryChan:
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
			}
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 2)

		}
	}
}

func (r *Raft) runCandidate() {
	timeoutCh := randomTimeout(time.Second * 2)
	for r.state == Candidate {
		select {
		case req := <-r.appendEntryChan:
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
			}
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 2)
			wg := sync.WaitGroup{}
			wg.Add(len(r.custer))
			for _, node := range r.custer {
				go func() {
					defer wg.Done()
					rsp, err := node.Vote(context.Background(), &api.VoteRequest{
						Term:         r.term,
						CandidateId:  r.id,
						LastLogIndex: 0,
						LastLogTerm:  0,
					})
					if err != nil {
						logrus.Error(err)
					}
					if rsp.VoteGranted {

					}
				}()
			}
			wg.Wait()

		}
	}
}

func (r *Raft) Quore() bool {

}
func randomTimeout(minVal time.Duration) <-chan time.Time {
	if minVal == 0 {
		return nil
	}
	extra := time.Duration(rand.Int63()) % minVal
	return time.After(minVal + extra)
}
