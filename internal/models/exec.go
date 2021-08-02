package models

import "encoding/json"

type ExecConfig struct {
	Enabled  bool
	Provider string
	Config   string
}

func (e *ExecConfig) DeserializeConfig(v interface{}) error {
	return json.Unmarshal([]byte(e.Config), v)
}

func (e *ExecConfig) SerializeConfig(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.Config = string(b)
	return nil
}

type ExecConfigRanna struct {
	Endpoint string
	Token    string
}

type ExecConfigJdoodle struct {
	ClientID     string
	ClientSecret string
}
