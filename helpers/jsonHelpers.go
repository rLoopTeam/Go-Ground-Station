package helpers

import (
	"rloop/Go-Ground-Station/gstypes"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

func DecodeNetworkingFile(path string) gstypes.Networking{
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}

	var config gstypes.Config
	json.Unmarshal(bytes, &config)
	fmt.Printf("parsed json config: %v \n",config)
	return config.Networking
}
