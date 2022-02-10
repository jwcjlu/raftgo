package raft

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jwcjlu/raftgo/api"
	"github.com/jwcjlu/raftgo/config"
	"io"
	"os"
)

type Log struct {
	file   *os.File
	data   []*api.LogEntry
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
	var readBytes int64
	for {
		entry, n, err := log.decodeEntry()
		if err != nil {
			if err == io.EOF {
			} else {
				if err = os.Truncate(filePath, readBytes); err != nil {
					return fmt.Errorf("raft.Log: Unable to recover: %v", err)
				}
			}
			break
		}
		readBytes += n
		log.data = append(log.data, entry)

	}
	return nil
}

func (log *Log) ApplyLogEntry(entry *api.LogEntry) error {
	err := log.encodeEntry(entry)
	if err != nil {
		return err
	}
	log.data = append(log.data, entry)
	return err
}

func (log *Log) encodeEntry(entry *api.LogEntry) error {
	data, err := proto.Marshal(entry)
	if err != nil {
		return err
	}
	err = WriteInt(log.file, len(data))
	if err != nil {
		return err
	}
	_, err = log.file.Write(data)
	return err
}

func (log *Log) decodeEntry() (*api.LogEntry, int64, error) {
	length, err := ReadInt(log.file)
	if err != nil {
		return nil, 0, err
	}
	if length == 0 {
		return nil, 0, nil
	}
	data := make([]byte, length)
	_, err = log.file.Read(data)
	if err != nil {
		return nil, 0, err
	}
	var entry api.LogEntry
	err = proto.Unmarshal(data, &entry)
	return &entry, int64(length + 8), err
}
