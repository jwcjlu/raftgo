package main

import (
	"github.com/jwcjlu/raftgo/api"
	"github.com/jwcjlu/raftgo/config"
	"github.com/jwcjlu/raftgo/raft"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"os"
	"sync"
)

func main() {
	//读取yaml配置文件
	file := os.Args[1]
	container := buildContainer(file)
	if err := container.Invoke(func(rpcServer *grpc.Server, service *raft.ServiceManager, apiService *api.OpenApiService) {
		listener, err := service.Start()
		if err != nil {
			panic(err)
		}
		service.RegisterService(rpcServer)
		// 4. 运行rpcServer，传入listener
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			_ = rpcServer.Serve(listener)
			wg.Done()
		}()
		go func() {
			apiService.Start()
			wg.Done()
		}()
		wg.Wait()
	}); err != nil {
		panic(err)
	}

}

func buildContainer(configFile string) *dig.Container {
	container := dig.New()
	var constructors []interface{}
	constructors = append(constructors, config.NewConf,
		grpc.NewServer, raft.NewServiceManager, raft.NewRaft, raft.NewService, api.NewOpenApiService)
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
