package raft

import (
	"fmt"
	"github.com/jwcjlu/raftgo/api"
	"github.com/jwcjlu/raftgo/config"
	"os"
)

type Log struct {
	file   *os.File
	data   []api.LogEntry
	logDir string
}

func (log *Log) Init(config *config.Config) error {
	log.logDir = config.Node.LogDir
	_, err := os.Stat(config.Node.LogDir)
	if err != nil {
		err = os.MkdirAll(config.Node.LogDir, os.ModePerm)

	}
	if err != nil {
		return err
	}
	file, err := os.Open(fmt.Sprintf("%s%s%s", config.Node.LogDir, "/", "raft.log"))
	if err != nil {
		return err
	}
	log.file = file
	return nil
}
