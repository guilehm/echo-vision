package filestorage

type FileStoragePort interface {
	GeneratePreSignedURL(fileKey FileKey, contentType string) (string, error)
}
