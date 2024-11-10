package manor

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	etcd "go.etcd.io/etcd/client/v3"
)

func Manorable(nodeStore *ast.NodeStore) {
	client, err := etcd.New(etcd.Config{
		Endpoints:   env.LighthouseConfig.Etcd.Endpoints,
		Username:    env.LighthouseConfig.Etcd.Username,
		Password:    env.LighthouseConfig.Etcd.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error().Msgf("Failed to connect to etcd: %v", err)
		return
	}
	defer client.Close()

	register(nodeStore)
}

func register(nodeStore *ast.NodeStore) error {
	jsonData, err := sonic.Marshal(nodeStore)
	if err != nil {
		return err
	}
	log.Info().Msgf("%v", string(jsonData))

	return nil
}
