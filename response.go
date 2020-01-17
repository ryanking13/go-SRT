package srt

import "encoding/json"

type responseParser struct {
	originalMsg []byte
	json        map[string]interface{}
	status      map[string]interface{}
	parsed      bool
}

// Parse parses message from SRT server
func (r *responseParser) Parse(msg []byte) error {
	err := json.Unmarshal(msg, &r.json)
	if err != nil {
		return err
	}

	resultMap := r.json["resultMap"].([]interface{})
	r.originalMsg = msg
	r.status = resultMap[0].(map[string]interface{})
	r.parsed = true

	return nil
}

// Success checks whether parsing is successed
func (r *responseParser) Success() bool {
	if r.status == nil {
		return false
	}
	if status := r.status["strResult"].(string); status == "SUCC" {
		return true
	} else if status == "FAIL" {
		return false
	} else {
		return false
	}
}

// Data is used to get parsed response data
func (r *responseParser) Data() map[string]interface{} {
	return r.json
}

func (r *responseParser) String() string {
	if !r.parsed {
		return ""
	}
	return string(r.originalMsg)
}
