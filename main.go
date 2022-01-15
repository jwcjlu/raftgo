package main

import (
	"github.com/raftgo/raft"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"os"
)

func main() {
	//读取yaml配置文件
	file := os.Args[1]
	container := buildContainer(file)


	if err := container.Invoke(func(rpcServer *grpc.Server, service *raft.ServiceManager) {
		listener, err := service.Start()
		if err != nil {
			panic(err)
		}
		service.RegisterService(rpcServer)
		// 4. 运行rpcServer，传入listener
		_ = rpcServer.Serve(listener)
	}); err != nil {
		panic(err)
	}

}

func buildContainer(config string) *dig.Container {
	container := dig.New()
	var constructors []interface{}
	constructors = append(constructors, raft.NewConf, raft.NewRaft, raft.NewServer,
		raft.NewServiceManager, grpc.NewServer)

	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			panic(err)
		}
	}
	container.Provide(func()string {
		return config
	})
	return container
}
