// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"backend/internal/repository"
	"backend/internal/server"
	"backend/pkg/app"
	"backend/pkg/log"
	"backend/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	db := repository.NewDB(viperViper, logger)
	sidSid := sid.NewSid()
	syncedEnforcer := repository.NewCasbinEnforcer(viperViper, logger, db)
	migrateServer := server.NewMigrateServer(db, logger, sidSid, syncedEnforcer)
	appApp := newApp(migrateServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewCasbinEnforcer)

var serverSet = wire.NewSet(server.NewMigrateServer)

// build App
func newApp(
	migrateServer *server.MigrateServer,
) *app.App {
	return app.NewApp(app.WithServer(migrateServer), app.WithName("demo-migrate"))
}
