//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/rganes5/go-gin-clean-arch/pkg/api"
	handler "github.com/rganes5/go-gin-clean-arch/pkg/api/handler"
	config "github.com/rganes5/go-gin-clean-arch/pkg/config"
	db "github.com/rganes5/go-gin-clean-arch/pkg/db"
	repository "github.com/rganes5/go-gin-clean-arch/pkg/repository"
	usecase "github.com/rganes5/go-gin-clean-arch/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase,
		repository.NewUserRepository, repository.NewAdminRepository,
		usecase.NewUserUseCase, usecase.NewAdminUseCase,
		handler.NewUserHandler, handler.NewAdminHandler,
		http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
