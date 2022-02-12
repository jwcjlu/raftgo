package raft

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jwcjlu/raftgo/config"
	proto2 "github.com/jwcjlu/raftgo/proto"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type Log struct {
	file   *os.File
	data   []*proto2.LogEntry
	logDir string
	temp   []*proto2.LogEntry
	mu     sync.Mutex
}

func (log *Log) reset() {
	log.temp = make([]*proto2.LogEntry, 0)
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
	logrus.Infof("data=%v", log.data)
	log.mu = sync.Mutex{}
	return nil
}

func (log *Log) ApplyLogEntry(entries []*proto2.LogEntry) error {
	for _, entry := range entries {
		err := log.encodeEntry(entry)
		if err != nil {
			return err
		}
		log.data = append(log.data, entry)
	}
	return nil
}
func (log *Log) commitLogEntry(leaderCommit int64) error {
	if len(log.temp) < 1 {
		return nil
	}
	logrus.Infof("commitLogEntry %v", log.temp)
	lastEntry := log.LastEntry()
	isTermDiff := lastEntry.CurrentTerm < log.temp[0].CurrentTerm
	if isTermDiff && log.temp[0].Index > 1 {
		return fmt.Errorf("data is error lastEntry=%v however temp[0]=%v", lastEntry, log.temp[0])
	}
	if !isTermDiff && log.temp[0].Index > lastEntry.Index {
		return fmt.Errorf("data is error lastEntry=%v however temp[0]=%v", lastEntry, log.temp[0])
	}
	nextIndex := int64(1)
	if !isTermDiff {
		nextIndex = lastEntry.Index
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	var appendEntries []*proto2.LogEntry
	var lastIndex int
	for index, entry := range log.temp {
		if entry.Index <= nextIndex && !isTermDiff {
			continue
		}
		if entry.Index > leaderCommit {
			break
		}
		lastIndex = index
		appendEntries = append(appendEntries, entry)
	}

	err := log.ApplyLogEntry(appendEntries)
	if err != nil {
		return err
	}

	log.temp = log.temp[lastIndex:]
	return nil
}
func (log *Log) temporaryLogEntry(entries []*proto2.LogEntry) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.temp = append(log.temp, entries...)
}

func (log *Log) encodeEntry(entry *proto2.LogEntry) error {
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
func (log *Log) NewAppendEntryRequest(raft *Raft, data *proto2.LogEntry, term int64, index int) *proto2.AppendEntriesRequest {
	entry := log.LastEntry()
	if index > 0 {
		entry = log.data[index-1]
	}
	if index == 0 {
		entry = &proto2.LogEntry{}
	}
	return &proto2.AppendEntriesRequest{
		Term:        term,
		PreLogIndex: entry.Index,
		PreLogTerm:  entry.CurrentTerm,
		Entry:       []*proto2.LogEntry{data},
	}
}

func (log *Log) decodeEntry() (*proto2.LogEntry, int64, error) {
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
	var entry proto2.LogEntry
	err = proto.Unmarshal(data, &entry)
	return &entry, int64(length + 9), err
}
func (log *Log) LastEntry() *proto2.LogEntry {
	if len(log.data) < 1 {
		return &proto2.LogEntry{}
	}
	return log.data[len(log.data)-1]
}

func (log *Log) DataLength() int {
	return len(log.data)
}
