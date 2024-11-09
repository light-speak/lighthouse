package handler

import (
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/manor"
)

func Manorable(nodeStore *ast.NodeStore) error {
	// Skip if not in cluster mode
	if env.LighthouseConfig.App.Mode != env.Cluster {
		return nil
	}

	// Initialize manor client
	err := manor.InitManorClient(env.LighthouseConfig.Manor.Host+":"+env.LighthouseConfig.Manor.Port, env.LighthouseConfig.App.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize manor client")
		return err
	}

	// Register node store with manor
	if err := manor.Register(nodeStore); err != nil {
		log.Error().Err(err).Msg("failed to register with manor")
		return err
	}

	log.Info().Msg("successfully registered with manor")
	return nil
}
