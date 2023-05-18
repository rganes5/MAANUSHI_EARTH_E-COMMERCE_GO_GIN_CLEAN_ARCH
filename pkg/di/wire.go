//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/rganes5/maanushi_earth_e-commerce/pkg/api"
	handler "github.com/rganes5/maanushi_earth_e-commerce/pkg/api/handler"
	config "github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
	db "github.com/rganes5/maanushi_earth_e-commerce/pkg/db"
	repository "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository"
	usecase "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase,
		repository.NewUserRepository, repository.NewAdminRepository, repository.NewProductRepository, repository.NewOtpRepository,
		usecase.NewUserUseCase, usecase.NewAdminUseCase, usecase.NewProductUseCase, usecase.NewOtpUseCase,
		handler.NewUserHandler, handler.NewAdminHandler, handler.NewProductHandler,
		http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
