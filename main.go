package main

import (
	"github.com/raftgo/config"
	"github.com/raftgo/openapi"
	"github.com/raftgo/raft"
	"github.com/raftgo/store"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"os"
	"sync"
)

func main() {
	//读取yaml配置文件
	file := os.Args[1]
	container := buildContainer(file)


	if err := container.Invoke(func(rpcServer *grpc.Server, service *raft.ServiceManager,apiService openapi.OpenApiService) {
		listener, err := service.Start()
		if err != nil {
			panic(err)
		}
		service.RegisterService(rpcServer)
		// 4. 运行rpcServer，传入listener
		wg:=sync.WaitGroup{}
		wg.Add(1)
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
	constructors = append(constructors, config.NewConf, raft.NewRaft, raft.NewServer,
		raft.NewServiceManager, grpc.NewServer,openapi.NewOpenApiService,store.NewMemoryStore)

	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			panic(err)
		}
	}
	container.Provide(func()string {
		return configFile
	})
	return container
}
