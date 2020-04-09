package core

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Config struct {
	Http   HttpConfig   `json:"http" yaml:"http"`
	Socket SocketConfig `json:"socket" yaml:"socket"`
	Traces  []TraceConfig `json:"traces" yaml:"traces"`
}

type HttpConfig struct {
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Domain  string        `json:"domain" yaml:"domain"`
	Headers []HttpHeader  `json:"headers"" yaml:"headers"`
}

type HttpHeader struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type SocketConfig struct {
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Ip      string        `json:"ip" yaml:"ip"`
	Port    int           `json:"port" yaml:"port"`
}

type TraceConfig struct {
	Id string				`json:"id" yaml:"id"`
	Name  string       `json:"name" yaml:"name"`
	Spans []SpanConfig `json:"spans" yaml:"spans"`
	Desc string `json:"desc" yaml:"desc"`

}

type SpanConfig struct {
	Id string				`json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Kind Kind   `json:"kind" yaml:"kind"`
	Desc string `json:"desc" yaml:"desc"`

	HttpEntry   HttpEntry   `json:"http_entry" yaml:"http_entry"`
	SocketEntry SocketEntry `json:"socket_entry" yaml:"socket_entry"`
}


type Kind int8

const (
	HTTP_KIND   Kind = 1
	SOCKET_KIND Kind = 2
)

type HttpMethod string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

type HttpEntry struct {
	Url     string            `json:"url" yaml:"url"`
	Method  HttpMethod        `json:"method" yaml:"method"`
	Params  []HttpParameter            `json:"params" yaml:"params"`
	Headers []HttpHeader `json:"headers" yaml:"headers"`
	State string `json:"state" yaml:"state"`
	Success []string `json:"success" yaml:"success"`
	Fail string `json:"fail" yaml:"fail"`
}

const (
	BREAK string = "break"
	CONTINUE string = "continue"
)


type HttpParameter struct {
	Name string		`json:"name" yaml:"name"`
	Value string	`json:"value" yaml:"value"`
}

type SocketEntry struct {
	Ip   string `json:"ip" yaml:"ip"`
	Port int    `json:"port" yaml:"port"`
}


func DefaultHttpConfig() HttpConfig {
	return HttpConfig{
		Timeout: 3,
	}
}

func DefaultSocketConfig() SocketConfig{
	return SocketConfig{
		Timeout: 3,
	}
}

func DefaultConfig() Config{
	return Config{
		Http:   DefaultHttpConfig(),
		Socket: DefaultSocketConfig(),
		Traces:  nil,
	}
}

func LoadConfig(path string) (Config, error) {
	c, err := ParseYAML(path)
	if err != nil {
		panic(err)
	}
	return c, nil
}

func ParseYAML(filename string) (Config, error) {
	c := DefaultConfig()
	yamlAbsPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(yamlAbsPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		panic(err)
	}
	return c, nil
}



