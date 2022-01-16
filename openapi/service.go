package openapi

import (
	"fmt"
	"github.com/raftgo/config"
	"github.com/raftgo/raft"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type OpenApiService struct {
	raft.LifeCycle
	config *config.Config
	raft   raft.Raft
}

func NewOpenApiService(config *config.Config, raft raft.Raft) OpenApiService {
	return OpenApiService{config: config, raft: raft}
}

func (api *OpenApiService) Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.PUT("put", api.put)
	e.GET("get/:key", api.get)
	e.Start(fmt.Sprintf("%s:%d", api.config.Node.Ip, api.config.Node.ApiPort))
}

type Entry struct {
	Key   string
	Value string
}

func (api *OpenApiService) put(c echo.Context) error {
	var entry Entry
	if err := c.Bind(&entry); err != nil {
		return err
	}
	logrus.Info("entry:", entry)
	err := api.raft.Put([]byte(entry.Key), []byte(entry.Value))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Hello, World!")
}
func (api *OpenApiService) get(c echo.Context) error {
	key := c.Param("key")
	logrus.Info("key:", key)
	value, err := api.raft.Get([]byte(key))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(value))
}
