package server

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(
		NewHealthServer,
		NewTransactionServer,
		NewUserServer,
		NewCategoryServer,
		NewAccountBookServer,
		NewFriendServer,
		NewAccountServer,
		NewChartServer,
	),
)
