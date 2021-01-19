package kafka

import (
	"errors"
	"log"
)

//实现消费者的接口
type MQInterface interface {
	OnMsgReceiver(Message)
	OnMsgError(error)
	OnMsgRebalance(Notification)
	OnClosed()
}

//消费者基类
type CBaseCustomer struct {
	c      Consumer
	isInit bool
}

func (this *CBaseCustomer) IsInit() bool {
	return this.isInit
}

func (this *CBaseCustomer) Start() error {
	if this.isInit {
		this.c.Start()
		return nil
	} else {
		return errors.New("请先初始化生产者")
	}
}

func (this *CBaseCustomer) Close() {
	if this.isInit {
		this.c.Stop()
		this.isInit = false
	}
	return
}

func (this *CBaseCustomer) Config(topic, consumerId string, v MQInterface) error {

	cfg := ConsumerConfig{
		UserName:      globalConfig.UserName,
		Password:      globalConfig.Password,
		Servers:       globalConfig.Servers,
		Topics:        []string{topic},
		ConsumerGroup: consumerId,
	}

	log.Println("准备初始化[生产者 ", consumerId, " ],这里耗时较长,等待一会")

	customer := Consumer{}
	if err := customer.Prepare(&cfg); err != nil {
		return errors.New("消费者初始化失败:" + err.Error())
	} else {
		this.c = customer
	}

	this.c.OnClosed = v.OnClosed
	this.c.OnMsgError = v.OnMsgError
	this.c.OnMsgRebalance = v.OnMsgRebalance
	this.c.OnMsgReceiver = v.OnMsgReceiver

	this.isInit = true

	if err := this.Start(); err != nil {
		return errors.New("消费者初始化失败:" + err.Error())
	}

	log.Println("初始化[生产者 ", consumerId, " ]成功", this.isInit)
	return nil
}
