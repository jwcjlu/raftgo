package raft

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"

	"github.com/jwcjlu/raftgo/api"
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
	Id    string
	Ip    string
	Port  int
	State StateEnum
	pool  Pool
}

func (node *Node) Init() {
	node.pool = Pool{factory: func() (io.Closer, error) {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", node.Ip, node.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}, resources: make(chan io.Closer, 3)}
}
func (node *Node) Vote(ctx context.Context, request *api.VoteRequest) (*api.VoteResponse, error) {
	conn, err := node.getConn()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer func() {
		cancel()
		if err != nil && conn != nil {
			conn.Close()
		}
		if conn != nil {
			node.pool.Release(conn)
		}
	}()
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
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer func() {
		cancel()
		if err != nil && conn != nil {
			conn.Close()
		}
		if conn != nil {
			node.pool.Release(conn)
		}
	}()
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

func (node *Node) Close() {
	node.pool.Close()
}
