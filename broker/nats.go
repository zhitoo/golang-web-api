package broker

import (
	"fmt"
	"log"

	"github.com/zhitoo/golang-web-api/config"

	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn

func InitNats(cfg *config.Config) error {
	var err error
	natsConn, err = nats.Connect(
		fmt.Sprintf("nats://%s:%s", cfg.Nats.Host, cfg.Nats.Port),
		nats.UserInfo("nats", cfg.Nats.Password),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetNats() *nats.Conn {
	return natsConn
}

func CloseNats() {
	if natsConn != nil {
		natsConn.Close()
	}
}
