package services

type AvDataService struct {
	ffw FileFlvWriter
	hfw HttpFlvWriter
}

func NewAvDataService() *AvDataService {
	return &AvDataService{}
}

func (avDataServie *AvDataService) notifyAll() {

}
