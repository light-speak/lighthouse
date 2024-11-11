package manor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	etcd "go.etcd.io/etcd/client/v3"
)

type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]*ServiceInfo
}

type ServiceInfo struct {
	Value     *RegisterValue
	NodeStore *ast.NodeStore
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*ServiceInfo),
	}
}

func (r *ServiceRegistry) Update(key string, value *RegisterValue) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Parse schema JSON into NodeStore
	nodeStore := &ast.NodeStore{
		Objects:    make(map[string]*ast.ObjectNode),
		Interfaces: make(map[string]*ast.InterfaceNode),
		Unions:     make(map[string]*ast.UnionNode),
		Enums:      make(map[string]*ast.EnumNode),
		Inputs:     make(map[string]*ast.InputObjectNode),
		Scalars:    make(map[string]*ast.ScalarNode),
	}

	if err := sonic.Unmarshal([]byte(value.Schema), nodeStore); err != nil {
		return fmt.Errorf("failed to unmarshal schema: %v", err)
	}

	r.services[key] = &ServiceInfo{
		Value:     value,
		NodeStore: nodeStore,
	}
	return nil
}

func (r *ServiceRegistry) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, key)
}

func (r *ServiceRegistry) GetAll() map[string]*ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*ServiceInfo, len(r.services))
	for k, v := range r.services {
		result[k] = v
	}
	return result
}

func (r *ServiceRegistry) GetService(key string) (*ServiceInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	service, exists := r.services[key]
	return service, exists
}

func Start() error {
	client, err := etcd.New(etcd.Config{
		Endpoints:   env.LighthouseConfig.Etcd.Endpoints,
		Username:    env.LighthouseConfig.Etcd.Username,
		Password:    env.LighthouseConfig.Etcd.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer client.Close()

	registry := NewServiceRegistry()
	if err := initServiceWatcher(client, registry, env.LighthouseConfig.App.Environment); err != nil {
		return err
	}

	// Print loaded schemas for debugging
	services := registry.GetAll()
	_, err = MergeSchemas(services)
	if err != nil {
		log.Error().Msgf("Failed to merge schemas: %v", err)
	}

	select {}
}

func initServiceWatcher(client *etcd.Client, registry *ServiceRegistry, environment env.AppEnvironment) error {
	servicePrefix := fmt.Sprintf("/manor/%s/", environment)

	// First get all existing services
	services, err := getServices(client, servicePrefix)
	if err != nil {
		return err
	}

	// Initialize registry with existing services
	for key, service := range services {
		if err := registry.Update(key, service); err != nil {
			log.Error().Msgf("Failed to update service %s: %v", key, err)
			continue
		}
	}
	log.Info().Msgf("Loaded %d existing services from %s", len(services), servicePrefix)

	// Start watching for changes
	go watchServices(client, registry, servicePrefix)

	return nil
}

func getServices(client *etcd.Client, prefix string) (map[string]*RegisterValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, prefix, etcd.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	services := make(map[string]*RegisterValue)
	for _, kv := range resp.Kvs {
		service := &RegisterValue{}
		if err := sonic.Unmarshal(kv.Value, service); err != nil {
			log.Warn().Msgf("Failed to unmarshal service data for key %s: %v", kv.Key, err)
			continue
		}
		key := string(kv.Key)
		services[key] = service
	}

	return services, nil
}

func watchServices(client *etcd.Client, registry *ServiceRegistry, prefix string) {
	watchChan := client.Watch(context.Background(), prefix, etcd.WithPrefix())

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			key := string(event.Kv.Key)

			switch event.Type {
			case etcd.EventTypePut:
				service := &RegisterValue{}
				if err := sonic.Unmarshal(event.Kv.Value, service); err != nil {
					log.Error().Msgf("Failed to unmarshal service update for key %s: %v", key, err)
					continue
				}
				if err := registry.Update(key, service); err != nil {
					log.Error().Msgf("Failed to update service %s: %v", key, err)
					continue
				}
				log.Info().Msgf("Service updated: %s", key)

				// Re-merge schemas after service update
				services := registry.GetAll()
				_, err := MergeSchemas(services)
				if err != nil {
					log.Error().Msgf("Failed to merge schemas after update: %v", err)
				}

			case etcd.EventTypeDelete:
				registry.Delete(key)
				log.Info().Msgf("Service deleted: %s", key)

				// Re-merge schemas after service deletion
				services := registry.GetAll()
				_, err := MergeSchemas(services)
				if err != nil {
					log.Error().Msgf("Failed to merge schemas after deletion: %v", err)
				}
			}
		}
	}
}
