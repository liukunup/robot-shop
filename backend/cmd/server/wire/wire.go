//go:build wireinject
// +build wireinject

package wire

import (
	"backend/internal/handler"
	"backend/internal/job"
	"backend/internal/repository"
	"backend/internal/server"
	"backend/internal/service"
	"backend/pkg/app"
	"backend/pkg/email"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"backend/pkg/server/http"
	"backend/pkg/sid"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewCasbinEnforcer,
	repository.NewUserRepository,
	repository.NewRoleRepository,
	repository.NewMenuRepository,
	repository.NewApiRepository,
	// more biz repository
	repository.NewRobotRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewRoleService,
	service.NewMenuService,
	service.NewApiService,
	// more biz service
	service.NewRobotService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewRoleHandler,
	handler.NewMenuHandler,
	handler.NewApiHandler,
	// more biz handler
	handler.NewRobotHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		email.NewEmail,
		newApp,
	))
}
