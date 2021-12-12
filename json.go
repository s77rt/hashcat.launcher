package hashcatlauncher

import (
	"encoding/json"
)

func MarshalJSON(v interface{}) []byte {
	json_raw, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	}
	return json_raw
}

func MarshalJSONS(v interface{}) string {
	// the S stands for string
	return string(MarshalJSON(v))
}
