package domain

type FileStorage interface {
	UploadFiles(path string, file []byte) error
	GetFile(path string) ([]byte, error)
}
