//go:build wireinject
// +build wireinject

package simple

import "github.com/google/wire"

func InitializedService(isErr bool) (*SimpleService, error) {
	wire.Build(NewSimpleRepository, NewSimpleService)

	return nil, nil
}

func InitializedDatabaseRepository() *DatabaseRepository {
	wire.Build(
		NewDatabasePostgreSQL, NewDatabaseMongoDB, NewDatabaseRepository,
	)

	return nil
}
