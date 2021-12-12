package raft

type LifeCycle interface {
	Init(conf *Config)

	Start()

	Stop()
}
