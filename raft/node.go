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
	nextIndex  int64
	matchIndex int64
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
	logrus.WithField("node", fmt.Sprintf("%s:%d", node.Ip, node.Port)).Info("startReplicate ...")
	index := node.raft.DataLength()
	negotiateFlag := false
	var err error
	for node.raft.state == Leader {

		if negotiateFlag {
			index++
		} else {
			if index > 0 {
				index--
			}
		}
		negotiateFlag, err = node.negotiate(index)
		if err != nil {
			logrus.Errorf("negotiate failure", err)
			return
		}
		if index == node.raft.DataLength()-1 && negotiateFlag {
			logrus.Info("startReplicate finish")
			return
		}
		if index == 0 && !negotiateFlag {
			logrus.Error("startReplicate failure not data")
			return
		}
	}
}

func (node *Node) negotiate(index int) (bool, error) {
	req := node.raft.NewAppendEntryRequest(index)
	req.IsApply = true
	rsp, err := node.AppendEntries(context.Background(), req)
	if err != nil {
		return false, err
	}
	return rsp.Success, nil
}
