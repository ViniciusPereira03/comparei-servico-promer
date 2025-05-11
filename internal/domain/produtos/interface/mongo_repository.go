package produtos_interface

type MongoRepository interface {
	SaveImage(image string, nome string) error
	GetImageByName(nome string) (string, error)
}
