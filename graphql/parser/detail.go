package parser

import (
	"encoding/json"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
)



func (p *Parser) NodeDetail(nodes map[string]ast.Node) {
	for _, node := range nodes {
		log.Info().Str("Description", node.GetDescription()).Str("Kind", node.GetKind().String()).Msg(node.GetName())
		switch node.GetKind() {
		case ast.KindObject:
			for _, field := range node.(*ast.ObjectNode).Fields {
				typeJson, err := json.Marshal(field.Type)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				argJson, err := json.Marshal(field.Args)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				directivesJson, err := json.Marshal(field.Directives)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				log.Warn().
					RawJSON("Type", typeJson).
					RawJSON("Args", argJson).
					RawJSON("Directives", directivesJson).
					Msg("\t" + field.Name)
			}
		case ast.KindInputObject:
			for _, field := range node.(*ast.InputObjectNode).Fields {
				typeJson, err := json.Marshal(field.Type)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				argJson, err := json.Marshal(field.Args)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				directivesJson, err := json.Marshal(field.Directives)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				log.Warn().
					RawJSON("Type", typeJson).
					RawJSON("Args", argJson).
					RawJSON("Directives", directivesJson).
					Msg("\t" + field.Name)
			}
		case ast.KindUnion:
			for _, typeName := range node.(*ast.UnionNode).TypeNames {
				log.Warn().Msgf("\tTypeName: %s", typeName)
			}
		case ast.KindEnum:
			for _, enumValue := range node.(*ast.EnumNode).EnumValues {
				directivesJson, err := json.Marshal(enumValue.Directives)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				log.Warn().RawJSON("Directives", directivesJson).Msgf("\tEnumValue: %s", enumValue.Name)
			}
		case ast.KindInterface:
			for _, field := range node.(*ast.InterfaceNode).Fields {
				typeJson, err := json.Marshal(field.Type)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				argJson, err := json.Marshal(field.Args)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				directivesJson, err := json.Marshal(field.Directives)
				if err != nil {
					log.Fatal().Err(err).Msg("")
				}
				log.Warn().
					RawJSON("Type", typeJson).
					RawJSON("Args", argJson).
					RawJSON("Directives", directivesJson).
					Msg(field.Name)
			}
		case ast.KindScalar:
			log.Warn().Msgf("\tScalar: %s", node.(*ast.ScalarNode).Name)
		}
	}
	for _, directive := range p.NodeStore.Directives {
		argJson, err := json.Marshal(directive.Args)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		log.Info().Str("Description", directive.Description).Str("Name", directive.Name).RawJSON("Args", argJson).Msg("")
	}
}