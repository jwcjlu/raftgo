package raft

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
}
