package raft

import (
	"context"
	"fmt"
	"github.com/jwcjlu/raftgo/proto"
	"google.golang.org/grpc/connectivity"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StateEnum int

const (
	Unknown StateEnum = iota
	Leader
	Follower
	Learner
	Candidate
)

type Node struct {
	Id         string
	Ip         string
	Port       int
	pool       Pool
	raft       *Raft
	term       int64
	nextIndex  int //当一个领导人刚获得权力的时候，他初始化所有的 nextIndex 值为自己的最后一条日志的 index 加 1
	matchIndex int
	isSync     bool
}

func (node *Node) Init(connCount int) {
	node.pool = Pool{factory: func() (io.Closer, error) {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Ip, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}, resources: make(chan io.Closer, connCount)}
}
func (node *Node) Vote(ctx context.Context, request *proto.VoteRequest) (*proto.VoteResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, deferFun := node.withDeferFunc(ctx, 500*time.Millisecond)
	defer deferFun(conn, err)
	rsp, err := proto.NewRaftServiceClient(conn).Vote(ctx, request)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return rsp, nil
}

func (node *Node) AppendEntries(ctx context.Context,
	request *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, deferFun := node.withDeferFunc(ctx, 500*time.Millisecond)
	defer deferFun(conn, err)
	rsp, err := proto.NewRaftServiceClient(conn).AppendEntries(ctx, request)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return rsp, nil
}
func (node *Node) Heartbeat(ctx context.Context,
	request *proto.HeartbeatRequest) (*proto.HeartbeatResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, deferFun := node.withDeferFunc(ctx, 500*time.Millisecond)
	defer deferFun(conn, err)
	rsp, err := proto.NewRaftServiceClient(conn).Heartbeat(ctx, request)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return rsp, nil
}
func (node *Node) getConn() (*grpc.ClientConn, error) {
	closer, err := node.pool.Acquire(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}
	conn, _ := closer.(*grpc.ClientConn)
	return conn, nil

}
func (node *Node) withDeferFunc(ctx context.Context, timeout time.Duration) (context.Context, func(*grpc.ClientConn, error)) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	return ctx, func(conn *grpc.ClientConn, err error) {
		cancel()
		if err != nil && conn != nil {
			conn.Close()
		}
		if conn != nil && conn.GetState() != connectivity.Shutdown {
			node.pool.Release(conn)
		}
	}

}

func (node *Node) Close() {
	node.pool.Close()
}

func (node *Node) startReplicate() {
	if node.isSync {
		return
	}
	node.isSync = true
	defer func() {
		node.isSync = false
	}()
	logrus.WithField("node", fmt.Sprintf("%s:%d", node.Ip, node.Port)).Info("startReplicate ...")
	negotiateFlag := false
	var err error
	index := 0
	for node.raft.state == Leader {
		if negotiateFlag {
			node.matchIndex--
		} else {
			node.matchIndex++
		}
		index = node.matchIndex
		if node.matchIndex < 0 {
			node.matchIndex = 0
			index = 0

		}
		if node.nextIndex == node.matchIndex+1 {
			logrus.Info("startReplicate finish")
			return
		}
		negotiateFlag, err = node.replicate(index)
		if err != nil {
			logrus.Errorf("negotiate failure", err)
			return
		}
		if index == 0 && !negotiateFlag {
			logrus.Error("startReplicate failure not data")
			return
		}
	}
}

func (node *Node) replicate(index int) (bool, error) {
	req := node.raft.NewAppendEntryRequest(index)
	req.IsApply = true
	rsp, err := node.AppendEntries(context.Background(), req)
	if err != nil {
		return false, err
	}
	return rsp.Success, nil
}

func (node *Node) reset(nextIndex int) {
	node.nextIndex = nextIndex
	node.matchIndex = nextIndex - 1

}
