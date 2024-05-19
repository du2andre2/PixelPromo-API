package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"pixelPromo/adapter/aws"
	"pixelPromo/adapter/config"
	"pixelPromo/domain/service"
	"pixelPromo/port/controller"
	"pixelPromo/port/route"
)

var AdapterModule = fx.Module("adapter",
	fx.Provide(
		aws.NewDynamoDb,
		aws.NewBucketS3,
		aws.NewConfigAWS,
		config.NewConfig,
		config.NewLogger,
	),
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		service.NewRepository,
	),
)

var PortModule = fx.Module("port",
	fx.Provide(
		route.NewServer,
		route.NewRoute,
		controller.NewController,
	),
)

var Module = fx.Options(
	AdapterModule,
	ServiceModule,
	PortModule,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	server route.Server,
) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

}
