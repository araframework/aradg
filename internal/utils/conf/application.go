package conf

import (
	"encoding/json"
	"os"
)

type Application struct {
	Hostname        string `json:"hostname"`
	ListenInterface string `json:"listen_interface"`
	PrivateKeyFile  string `json:"private_key_file"`
	PublicKeyFile   string `json:"public_key_file"`
	LogFile         string `json:"log_file"`
}

// load conf
func load() (*Application, error) {

	// load conf
	file, err := os.Open("conf/application.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	app := &Application{}
	err = decoder.Decode(app)
	// TODO check all fields valid
	// TODO check not empty, valid interface, valid key file, valid log file path etc.
	return app, err
}
