package Service

import (
	"io"
	"os"
	"time"
	"encoding/json"
)

type Config struct {
	Req      []Request     `json:"requests"`
	Host     string        `json:"host"`
	Duration time.Duration `json:"duration"`
}

type Request struct {
	Method      string        `json:"method"`
	Endpoint    string        `json:"endpoint"`
	Data        string        `json:"data"`
	Header      string        `json:"header"`
	Connections int           `json:"connections"`
	Rate        time.Duration `json:"rate"`
	ConcurrencyLimit int 	  `json:"concurrencyLimit"`
}

func FromJSON(path string) (*Config, error) {
	f, err := os.Open(path);
	if err != nil {
		return nil, err;
	}
	defer f.Close();

	data, err := io.ReadAll(f);
	if err != nil {
		return nil, err;
	}

	var conf Config;

	err = json.Unmarshal(data, &conf);
	if err != nil {
		return nil, err;
	}

	return &conf, nil;
}