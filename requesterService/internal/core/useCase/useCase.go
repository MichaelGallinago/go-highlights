package useCase

type UseCase struct {
	memeClient      MemeClient
	requesterServer RequesterServer
}

func NewUseCase(memeClient MemeClient, requesterServer RequesterServer) UseCase {
	return UseCase{
		memeClient:      memeClient,
		requesterServer: requesterServer,
	}
}
