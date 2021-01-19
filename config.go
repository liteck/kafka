package kafka

import "time"

type ConsumerConfig struct {
	Topics        []string
	Servers       []string
	UserName      string
	Password      string
	ConsumerGroup string
}

type ProducerConfig struct {
	Servers  []string
	Ak       string
	Password string
}

var globalConfig Config

type Config struct {
	Servers  []string
	UserName string
	Password string
	Cert     []byte
}

type Message struct {
	Key   string
	Msg   string
	Topic string
	Time  time.Time
}

type Notification struct {
	Claimed  map[string][]int32
	Current  map[string][]int32
	Released map[string][]int32
}
