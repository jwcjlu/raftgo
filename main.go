package main

import (
	"github.com/jwcjlu/raftgo/config"
	"github.com/jwcjlu/raftgo/raft"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"os"
)

func main() {
	//读取yaml配置文件
	file := os.Args[1]
	container := buildContainer(file)
	if err := container.Invoke(func(rpcServer *grpc.Server, service *raft.ServiceManager) {
		_, err := service.Start()
		if err != nil {
			panic(err)
		}
		service.RegisterService(rpcServer)
		// 4. 运行rpcServer，传入listener
	}); err != nil {
		panic(err)
	}

}

func buildContainer(configFile string) *dig.Container {
	container := dig.New()
	var constructors []interface{}
	constructors = append(constructors, config.NewConf,
		grpc.NewServer, raft.NewServiceManager)
	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			panic(err)
		}
	}
	container.Provide(func() string {
		return configFile
	})
	return container
}
