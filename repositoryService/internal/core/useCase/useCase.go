package useCase

import "repositoryService/internal/lib/postgresclient"

type UseCase struct {
	searchServer  SearchServer
	publishServer PublishServer
	db            DB
}

func NewUseCase(searchServer SearchServer, publishServer PublishServer) UseCase {
	return UseCase{
		searchServer:  searchServer,
		publishServer: publishServer,
		db:            &postgresclient.Instance,
	}
}

func (uc UseCase) CloseDB() {
	uc.db.Close()
}
