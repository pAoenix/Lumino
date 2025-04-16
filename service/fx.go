package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewUserService,
			fx.As(new(UserIconDownloader))),
		fx.Annotate(
			NewCategoryService,
			fx.As(new(CategoryDownloader))),
		NewTransactionService,
		NewUserService,
		NewCategoryService,
		NewAccountBookService,
		NewFriendService,
		NewAccountService,
	),
)
