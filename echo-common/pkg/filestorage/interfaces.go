package filestorage

type FileStoragePort interface {
	GeneratePreSignedURL(fileKey FileKey, contentType string) (string, error)
	GenerateFileURL(fileKey FileKey) (string, error)
}
