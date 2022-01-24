package raft

import (
	"context"
	"fmt"
	"google.golang.org/grpc/connectivity"
	"io"
	"time"

	"github.com/jwcjlu/raftgo/api"

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
func (node *Node) Vote(ctx context.Context, request *api.VoteRequest) (*api.VoteResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, deferFun := node.withDeferFunc(ctx, 500*time.Millisecond)
	defer deferFun(conn, err)
	rsp, err := api.NewRaftServiceClient(conn).Vote(ctx, request)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return rsp, nil
}

func (node *Node) AppendEntries(ctx context.Context,
	request *api.AppendEntriesRequest) (*api.AppendEntriesResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, deferFun := node.withDeferFunc(ctx, 500*time.Millisecond)
	defer deferFun(conn, err)
	rsp, err := api.NewRaftServiceClient(conn).AppendEntries(ctx, request)
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

}
