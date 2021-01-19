package kafka

import (
	"encoding/json"
	"errors"
	"log"
)

//实例化一个生产者
type PProducer struct {
	p      Producer
	inited bool
}

var globalProducer PProducer

func GetProducer() PProducer {
	return globalProducer
}

func (this *PProducer) Init() error {
	log.Println("[生产者模块] 初始化")
	globalProducer = PProducer{}

	cfg := ProducerConfig{
		Servers:  globalConfig.Servers,
		Ak:       globalConfig.UserName,
		Password: globalConfig.Password,
	}

	this.p = Producer{}
	this.inited = false
	if err := this.p.Prepare(&cfg); err != nil {
		this.p = Producer{}
		return errors.New("消息队列生产者初始化失败:" + err.Error())
	} else {
		this.inited = true
	}

	log.Println("[生产者模块] 初始化成功")
	return nil
}

func (this PProducer) Send(v interface{}, topic, key string) error {
	//把当前的数据重新发送一次
	if _x, err := json.Marshal(v); err != nil {
		return err
	} else {
		return this.SendString(string(_x), topic, key)
	}
}

func (this PProducer) SendString(s, topic, key string) error {
	return this.send(topic, s, key)
}

func (this PProducer) send(topic, content, key string) error {
	if this.inited == false {
		return errors.New("消息队列生产者未初始化")
	}
	//把当前的数据重新发送一次
	//uuid.Must(uuid.NewV4()).String()
	if err := this.p.SendMsg(topic, key, content); err != nil {
		return err
	} else {
		return nil
	}
}
