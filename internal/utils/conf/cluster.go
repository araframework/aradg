package conf

import (
	"encoding/json"
	"os"
)

type Cluster struct {
	Hostname        string `json:"hostname"`
	ListenInterface string `json:"listen_interface"`
	PrivateKeyFile  string `json:"private_key_file"`
	PublicKeyFile   string `json:"public_key_file"`
	LogFile         string `json:"log_file"`
}

// load conf
func loadClusterConf() (*Cluster, error) {

	// load conf
	file, err := os.Open("conf/cluster.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cluster := &Cluster{}
	err = decoder.Decode(cluster)
	// TODO check all fields valid
	// TODO check not empty, valid interface, valid key file, valid log file path etc.
	return cluster, err
}
