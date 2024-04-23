package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"pixelPromo/adapter/aws"
	"pixelPromo/domain/service"
	"pixelPromo/port/daemon"
	"pixelPromo/port/route"
)

var AdapterModule = fx.Module("adapter",
	fx.Provide(
		aws.NewDynamoDb,
		aws.NewBucketS3,
		aws.NewConfigAWS,
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
		daemon.NewDaemon,
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
	daemon daemon.Daemon,
) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go daemon.Run()
			go server.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			daemon.Stop()
			return nil
		},
	})

}
