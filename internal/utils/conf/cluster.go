package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClusterOption struct {
	Network Network `json:"network"`
}

type Network struct {
	Interface string              `json:"interface"`
	Join      map[string][]string `json:"join"`
}

var option *ClusterOption

// load conf
func LoadCluster() *ClusterOption {
	if option != nil {
		return option
	}

	file, err := os.Open("conf/cluster.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	option = &ClusterOption{}
	err = decoder.Decode(option)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// TODO check all fields valid
	return option
}
