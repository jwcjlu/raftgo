package raft

import "github.com/jwcjlu/raftgo/config"

type LifeCycle interface {
	Init(conf *config.Config)

	Start()

	Stop()
}
