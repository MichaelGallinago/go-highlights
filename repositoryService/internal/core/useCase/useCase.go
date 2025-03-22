package useCase

type UseCase struct {
	searchServer  SearchServer
	publishServer PublishServer
	db            DB
}

func NewUseCase(searchServer SearchServer, publishServer PublishServer, db DB) UseCase {
	return UseCase{
		searchServer:  searchServer,
		publishServer: publishServer,
		db:            db,
	}
}

func (uc UseCase) CloseDB() {
	uc.db.Close()
}
