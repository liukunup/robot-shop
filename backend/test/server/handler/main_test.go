package handler

import (
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/pkg/config"
	jwt_ "backend/pkg/jwt"
	"backend/pkg/log"
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var (
	adminUserID uint = 1
)
var logger *log.Logger
var hdl *handler.Handler
var jwt *jwt_.JWT
var router *gin.Engine

func TestMain(m *testing.M) {
	fmt.Println("begin")
	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		fmt.Println("Setenv error", err)
	}
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	// modify log directory
	logPath := filepath.Join("../../../", conf.GetString("log.log_file_name"))
	conf.Set("log.log_file_name", logPath)

	logger = log.NewLog(conf)
	hdl = handler.NewHandler(logger)

	// 创建一个mock的TokenStore
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	mockTokenStore := mock_repository.NewMockTokenStore(ctrl)

	// 正确初始化JWT
	jwt = jwt_.NewJwt(conf, mockTokenStore)
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func performRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func generateAccessToken(t *testing.T) string {
	ctx := context.Background()
	tokenPair, err := jwt.GenerateTokenPair(ctx, adminUserID, "")
	if err != nil {
		t.Error(err)
		return ""
	}
	return tokenPair.AccessToken
}

func newHttpExcept(t *testing.T, router *gin.Engine) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			// httpexpect.NewDebugPrinter(t, true),
		},
	})
}
