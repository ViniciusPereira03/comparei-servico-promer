package produtos_interface

type MongoRepository interface {
	SaveImage(image string, barcode string) error
	GetImageByBarcode(barcode string) ([]byte, string, error)
}
