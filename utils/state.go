package utils

import (
	"os"
	"encoding/json"
	types "htty/types"
)

//loads htty state from file (json type) to get a httystate object
func LoadState(filepath string) (types.HttyState) {
	var fileBuf []byte
	var err error
	fileBuf, err = os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var stateBuf types.HttyState
	err = json.Unmarshal(fileBuf, &stateBuf)
	if err != nil {
		panic(err)
	}
	return stateBuf
}


//take current httystate from req/res & save into file
func SaveState(curstate types.HttyState, filepath string) error {
	stateJson, err :=  json.Marshal(curstate);
	if err != nil { return err }
	err = WriteFileContents(filepath, string(stateJson))
	if err != nil { return err }
	return nil
}
