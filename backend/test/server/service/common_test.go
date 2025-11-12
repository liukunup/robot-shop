package service_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"backend/pkg/config"
	"backend/pkg/email"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/sid"
)

var (
	logger *log.Logger
	j      *jwt.JWT
	sf     *sid.Sid
	em     *email.Email
)

func TestMain(m *testing.M) {
	fmt.Println("begin service tests")

	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	j = jwt.NewJwt(conf, nil)
	sf = sid.NewSid()
	em = email.NewEmail(conf)

	code := m.Run()
	fmt.Println("service tests end")

	os.Exit(code)
}
