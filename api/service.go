package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jwcjlu/raftgo/config"
	"github.com/jwcjlu/raftgo/raft"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type OpenApiService struct {
	raft.LifeCycle
	config *config.Config
	raft   *raft.Raft
}

func NewOpenApiService(config *config.Config, raft *raft.Raft) *OpenApiService {
	return &OpenApiService{config: config, raft: raft}
}

func (api *OpenApiService) Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.PUT("put", api.put)
	//e.GET("get/:key", api.get)
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
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	api.raft.CmdCh <- data
	result := <-api.raft.CmdRsp
	if !result {
		return fmt.Errorf("failure")
	}
	return c.String(http.StatusOK, "success!")
}

/*func (api *OpenApiService) get(c echo.Context) error {
	key := c.Param("key")
	logrus.Info("key:", key)
	value, err := api.raft.Get([]byte(key))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(value))
}
*/
