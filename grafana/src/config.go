package main

import (
	"fmt"
)

func constructDataSourceConfig(peers map[string]interface{}) ConfigDataSources {
	config := ConfigDataSources{}

	for identity, rawPeer := range peers {
		fmt.Printf("***** [constructConfig] PEER: %#v\n", rawPeer)
		peer := rawPeer.(map[string]interface{})

		targetURL := fmt.Sprintf("%s://%s:%d", peer["protocol"], peer["address"], int(peer["port"].(float64)))
		dataSource := ConfigDataSource{
			Access:    "proxy",
			Editable:  false,
			IsDefault: true,
			OrigID:    1,
			Type:      "prometheus",
			Name:      "berliozprom" + identity,
			URL:       targetURL,
			Version:   1,
		}

		config.DataSources = append(config.DataSources, &dataSource)
	}

	return config
}
