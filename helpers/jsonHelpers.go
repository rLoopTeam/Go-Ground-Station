package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"rloop/Go-Ground-Station-1/gstypes"
)

func DecodeNetworkingFile(path string) (gstypes.Networking, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	//TODO: Check for missing network
	var config gstypes.Config
	json.Unmarshal(bytes, &config)
	fmt.Printf("parsed json config: %v \n", config)
	return config.Networking, err
}
