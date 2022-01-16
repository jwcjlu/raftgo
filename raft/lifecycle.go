package raft

import "github.com/raftgo/config"

type LifeCycle interface {
	Init(conf *config.Config)

	Start()

	Stop()
}
