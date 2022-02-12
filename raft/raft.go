package raft

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jwcjlu/raftgo/config"
	"github.com/jwcjlu/raftgo/proto"

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
	term                    int64
	isLeader                bool
	lastApplied             int64
	commitIndex             int64
	lastTerm                int64
	lastIndex               int64
	index                   int64
	custer                  []*Node
	state                   StateEnum
	appendEntryChan         chan *proto.AppendEntriesRequest
	appendEntryResponseChan chan *proto.AppendEntriesResponse
	heartbeatChan           chan *proto.HeartbeatRequest
	heartbeatResponseChan   chan *proto.HeartbeatResponse
	LifeCycle
	mu     sync.RWMutex
	CmdCh  chan []byte
	CmdRsp chan bool
	log    *Log
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
	r.appendEntryResponseChan = make(chan *proto.AppendEntriesResponse, 1)
	r.appendEntryChan = make(chan *proto.AppendEntriesRequest, 1)
	r.heartbeatResponseChan = make(chan *proto.HeartbeatResponse, 1)
	r.heartbeatChan = make(chan *proto.HeartbeatRequest, 1)
	r.CmdCh = make(chan []byte, 1)
	r.CmdRsp = make(chan bool, 1)
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
	entry := r.log.LastEntry()
	r.term = entry.CurrentTerm
	r.lastApplied = entry.Index
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
	go r.startLeader()
	r.log.reset()
	for r.state == Leader {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Info("runLeader")
		select {
		case req := <-r.appendEntryChan:
			rsp := &proto.AppendEntriesResponse{Term: r.term, Success: false}
			if req.Term > r.term {
				r.state = Follower
				r.isLeader = false
				rsp.Success = true
			}
			r.appendEntryResponseChan <- rsp
		case req := <-r.heartbeatChan:
			rsp := &proto.HeartbeatResponse{Term: r.term, Success: false}
			if req.Term >= r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
				rsp.Success = true
			}
			r.heartbeatResponseChan <- rsp
		case <-timeoutCh:
			for _, n := range r.custer {
				go func(node *Node) {
					rsp, err := node.Heartbeat(context.Background(), &proto.HeartbeatRequest{
						Term:         r.term,
						LeaderId:     r.leaderId,
						LeaderCommit: r.commitIndex,
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

func (r *Raft) startLeader() {
	close(r.CmdCh)
	r.CmdCh = make(chan []byte, 1)
	r.index = 0
	r.commitIndex = 0
	logrus.Info("startLeader....")
	for r.state == Leader {
		cmd := <-r.CmdCh
		entry := &proto.LogEntry{
			CurrentTerm: r.term,
			Index:       r.index + 1,
			Data:        cmd,
		}
		req := r.log.NewAppendEntryRequest(r, entry, r.term, -1)
		voteGranted := 1
		wg := sync.WaitGroup{}
		wg.Add(len(r.custer))
		for _, n := range r.custer {
			go func(node *Node) {
				defer wg.Done()
				rsp, err := node.AppendEntries(context.Background(), req)
				if err != nil {
					logrus.Error(err)
					return
				}
				if rsp.Success {
					voteGranted++
				} else {
					go n.startReplicate()
				}
			}(n)
		}
		wg.Wait()
		flag := false
		if voteGranted >= r.QuorumSize() {
			r.log.ApplyLogEntry([]*proto.LogEntry{entry})
			r.commitIndex = entry.Index
			r.index = entry.Index
			flag = true
		}
		r.CmdRsp <- flag
	}

}

/*
1、如果term < currentTerm返回 false
2、如果 votedFor 为空或者为 candidateId，并且候选人的日志至少和自己一样新，那么就投票给他
*/
func (r *Raft) HandlerVote(request *proto.VoteRequest) (*proto.VoteResponse, error) {
	rsp := proto.VoteResponse{Term: r.term, VoteGranted: false}
	if request.Term < r.term && r.compareLogEntry(request) {
		rsp.VoteGranted = true
	}
	return &rsp, nil
}

func (r *Raft) compareLogEntry(request *proto.VoteRequest) bool {
	entry := r.log.LastEntry()
	if entry.Index <= request.LastLogIndex && entry.CurrentTerm <= request.Term {
		return true
	}
	return false
}

func (r *Raft) runFollower() {
	timeoutCh := randomTimeout(time.Second * 3)
	lastTime := time.Now()
	for r.state == Follower {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Info("runFollower")
		select {
		case req := <-r.appendEntryChan:
			lastTime = time.Now()
			r.HandlerFollowerAppendEntry(req)
		case <-timeoutCh:
			timeoutCh = randomTimeout(time.Second * 3)
			if time.Now().Sub(lastTime).Seconds() < 3*time.Second.Seconds() {
				continue
			}
			r.state = Candidate
		case req := <-r.heartbeatChan:
			timeoutCh = randomTimeout(time.Second * 3)
			rsp := &proto.HeartbeatResponse{Term: r.term, Success: false}
			if req.Term >= r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
				rsp.Success = true
			}
			r.log.commitLogEntry(req.LeaderCommit)
			r.heartbeatResponseChan <- rsp

		}
	}
}

/*
1、返回假 如果领导人的任期小于接收者的当前任期（译者注：这里的接收者是指跟随者或者候选人）
2、返回假 如果接收者日志中没有包含这样一个条目 即该条目的任期在 prevLogIndex 上能和 prevLogTerm 匹配上
（注：在接收者日志中 如果能找到一个和 prevLogIndex 以及 prevLogTerm 一样的索引和任期的日志条目 则继续执行下面的步骤 否则返回假
3、如果一个已经存在的条目和新条目（注：即刚刚接收到的日志条目）发生了冲突（因为索引相同，任期不同），
那么就删除这个已经存在的条目以及它之后的所有条目
4、追加日志中尚未存在的任何新条目
5、如果领导人的已知已提交的最高日志条目的索引大于接收者的已知已提交最高日志条目的索引（leaderCommit > commitIndex），
则把接收者的已知已经提交的最高的日志条目的索引commitIndex 重置为 领导人的已知已经提交的最高的日志条目的索引
leaderCommit 或者是 上一个新条目的索引 取两者的最小值
*/
func (r *Raft) HandlerFollowerAppendEntry(req *proto.AppendEntriesRequest) {
	rsp := &proto.AppendEntriesResponse{Term: r.term, Success: false}
	if req.Term < r.term {
		r.appendEntryResponseChan <- rsp
		return
	}
	logrus.Infof("AppendEntriesRequest:%v", req)
	entry := r.log.LastEntry()
	if entry.Index == req.PreLogIndex && entry.CurrentTerm == req.PreLogTerm {
		rsp.Success = true
	}
	if req.IsApply {
		r.log.ApplyLogEntry(req.Entry)
	} else if rsp.Success {
		r.log.temporaryLogEntry(req.Entry)
	}
	r.appendEntryResponseChan <- rsp
}

func (r *Raft) runCandidate() {
	timeoutCh := randomTimeout(time.Second * 3)
	isVote := true
	for r.state == Candidate {
		logrus.WithField("node", fmt.Sprintf("%s:%d", r.ip, r.port)).Infof("runCandidate term=%d", r.term)
		select {
		case req := <-r.heartbeatChan:
			timeoutCh = randomTimeout(time.Second * 3)
			rsp := &proto.HeartbeatResponse{Term: r.term, Success: false}
			if req.Term >= r.term {
				r.state = Follower
				r.isLeader = false
				r.leaderId = req.LeaderId
				rsp.Success = true
			}
			r.heartbeatResponseChan <- rsp
		case <-timeoutCh:
			if !isVote {
				timeoutCh = randomTimeout(time.Second * 3)
				continue
			}
			isVote = false
			voteGranted := 1
			r.term++
			wg := sync.WaitGroup{}
			entry := r.log.LastEntry()
			wg.Add(len(r.custer))
			for _, n := range r.custer {
				go func(node *Node, term int64) {
					defer wg.Done()
					rsp, err := node.Vote(context.Background(), &proto.VoteRequest{
						Term:         term,
						CandidateId:  r.id,
						LastLogIndex: entry.Index,
						LastLogTerm:  entry.CurrentTerm,
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
func (r *Raft) DataLength() int {
	return r.log.DataLength()
}
func (r *Raft) IndexEntry(index int) *proto.LogEntry {
	return r.log.data[index]
}

func (r *Raft) NewAppendEntryRequest(index int) *proto.AppendEntriesRequest {
	data := r.log.data[index]
	return r.log.NewAppendEntryRequest(r, data, data.CurrentTerm, index)
}
func randomTimeout(minVal time.Duration) <-chan time.Time {
	if minVal == 0 {
		return nil
	}
	extra := time.Duration(rand.Int63()) % minVal
	return time.After(minVal + extra)
}
