package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/http"
	"pixelPromo/adapter/repository"
	"pixelPromo/adapter/storage"
	"pixelPromo/config"
	"pixelPromo/domain/service"
)

var AdapterModule = fx.Module("adapter",
	fx.Provide(
		aws.NewConfigAWS,
		storage.NewBucketS3Storage,
		repository.NewDynamoDBRepository,
		http.NewRouter,
		http.NewController,
	),
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		service.NewService,
	),
)

var ConfigModule = fx.Module("config",
	fx.Provide(
		config.NewConfig,
		config.NewLogger,
	),
)

var Module = fx.Options(
	AdapterModule,
	ServiceModule,
	ConfigModule,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	router http.Router,
) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go router.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

}
