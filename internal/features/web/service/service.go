package web_service

type service struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) ([]byte, error)
}

func NewService(webRepository WebRepository) *service {
	return &service{
		webRepository: webRepository,
	}
}
