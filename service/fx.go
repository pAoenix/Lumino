package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewTransactionService,
		NewUserService,
		NewCategoryService,
		NewAccountBookService,
	),
)
