package elasticsearch

import (
	"bytes"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/light-speak/lighthouse/log"
)

type ElasticsearchOutput struct {
	client *elasticsearch.Client
	index  string
}

var client *elasticsearch.Client
var index string

func InitElasticsearchClient(Enable bool, Host string, Port string, User string, Password string, Index string) error {
	if !Enable {
		return nil
	}

	esConfig := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", Host, Port),
		},
		Username: User,
		Password: Password,
	}

	var err error
	client, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return fmt.Errorf("error creating the elasticsearch client: %s", err)
	}

	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("error connecting to elasticsearch: %s", err)
	}

	log.Info().Msgf("successfully connected to elasticsearch. cluster info: %s", res)

	return nil
}

func NewElasticsearchOutput() (*ElasticsearchOutput, error) {
	return &ElasticsearchOutput{
		client: client,
		index:  index,
	}, nil
}

func (e *ElasticsearchOutput) Write(p []byte) (n int, err error) {
	_, err = e.client.Index(e.index, bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
