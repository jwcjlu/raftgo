package raft

import (
	"fmt"
	"github.com/golang/protobuf/proto"
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
	filePath := fmt.Sprintf("%s%s%s", config.Node.LogDir, "/", "raft.log")
	_, err := os.Stat(filePath)
	if err != nil {
		err = os.MkdirAll(log.logDir, os.ModePerm)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		log.file = file
		return nil
	}
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		return err
	}
	log.file = file
	return nil
}

func (log *Log) ApplyLogEntry(entry *api.LogEntry) error {
	data, err := proto.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = log.file.Write(data)
	return err
}
