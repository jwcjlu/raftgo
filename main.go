package main

import (


	"go.uber.org/dig"
	"google.golang.org/grpc"
	"github.com/raftgo/raft"
)

func main() {
	container := buildContainer()
	//读取yaml配置文件

	if err := container.Invoke(func(rpcServer *grpc.Server, service *raft.ServiceManager) {
		listener, err := service.Start()
		if err!=nil{
			panic(err)
		}
		// 4. 运行rpcServer，传入listener
		_ = rpcServer.Serve(listener)
	}); err != nil {
		panic(err)
	}

}

func buildContainer() *dig.Container {
	container := dig.New()
	var constructors []interface{}
	constructors = append(constructors, raft.NewConf, raft.NewRaft, raft.NewServer, grpc.NewServer)

	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			panic(err)
		}
	}
	return container
}
