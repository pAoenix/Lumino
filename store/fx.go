package store

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewPgDB,
		NewTransactionStore,
		NewUserStore,
		NewCategoryStore,
	),
)
