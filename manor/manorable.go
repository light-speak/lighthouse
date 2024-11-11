package manor

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"

	etcd "go.etcd.io/etcd/client/v3"
)

type RegisterStatus string

const (
	RegisterStatusOnline  RegisterStatus = "online"
	RegisterStatusOffline RegisterStatus = "offline"
)

type RegisterValue struct {
	Schema      string         `json:"schema"`
	Address     string         `json:"address"`
	Status      RegisterStatus `json:"status"`
	Version     int            `json:"version"`
	LastUpdated int64          `json:"last_updated"`
}

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := register(client, nodeStore, ctx); err != nil {
		log.Error().Msgf("Failed to register schema to etcd: %v", err)
	}
	select {}
}

func register(client *etcd.Client, nodeStore *ast.NodeStore, ctx context.Context) error {
	// Filter out nodes starting with __
	filteredStore := &ast.NodeStore{
		Objects:    make(map[string]*ast.ObjectNode),
		Interfaces: make(map[string]*ast.InterfaceNode),
		Unions:     make(map[string]*ast.UnionNode),
		Enums:      make(map[string]*ast.EnumNode),
		Inputs:     make(map[string]*ast.InputObjectNode),
		Scalars:    make(map[string]*ast.ScalarNode),
	}

	// Copy objects excluding __ prefixed ones
	for name, obj := range nodeStore.Objects {
		filteredStore.Objects[name] = obj
	}

	// Copy other node types excluding __ prefixed ones
	for name, node := range nodeStore.Interfaces {
		filteredStore.Interfaces[name] = node
	}

	for name, node := range nodeStore.Unions {
		filteredStore.Unions[name] = node
	}

	for name, node := range nodeStore.Enums {
		filteredStore.Enums[name] = node
	}

	for name, node := range nodeStore.Inputs {
		filteredStore.Inputs[name] = node
	}

	for name, node := range nodeStore.Scalars {
		filteredStore.Scalars[name] = node
	}

	jsonData, err := sonic.Marshal(filteredStore)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}
	serviceName := utils.SnakeCase(env.LighthouseConfig.App.Name)

	value := RegisterValue{
		Schema:      string(jsonData),
		Address:     env.LighthouseConfig.Server.Endpoint,
		Status:      RegisterStatus(env.LighthouseConfig.App.Environment),
		Version:     env.LighthouseConfig.App.Version,
		LastUpdated: time.Now().Unix(),
	}

	jsonValue, err := sonic.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}

	// Create ETCD key with service name and version
	key := fmt.Sprintf("/manor/%s/%s",
		env.LighthouseConfig.App.Environment,
		serviceName,
	)

	// Put schema to ETCD with lease
	lease, err := client.Grant(ctx, 30) // 30 seconds TTL
	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}

	_, err = client.Put(ctx, key, string(jsonValue), etcd.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to put schema to etcd: %v", err)
	}

	// Keep lease alive
	keepAliveCh, err := client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("failed to keep lease alive: %v", err)
	}

	// Handle keep alive responses in background
	go func() {
		for {
			select {
			case resp, ok := <-keepAliveCh:
				if !ok || resp == nil {
					log.Error().Msg("Lost etcd lease keep alive")
					return
				}
			case <-ctx.Done():
				log.Info().Msg("Context cancelled, stopping etcd lease keep alive")
				return
			}
		}
	}()

	log.Info().Msgf("Successfully registered schema to etcd: %s", key)
	return nil
}
