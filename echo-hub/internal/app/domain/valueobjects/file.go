package valueobjects

type File struct {
	Filepath    string `json:"filepath"`
	Filename    string `json:"filename"`
	Filesize    int64  `json:"filesize"`
	ContentType string `json:"contentType"`
}

func NewFile(filepath, filename, contentType string, filesize int64) *File {
	return &File{
		Filepath:    filepath,
		Filename:    filename,
		Filesize:    filesize,
		ContentType: contentType,
	}
}

func (f *File) IsValid() bool {
	return f.Filepath != "" && f.Filename != "" && f.Filesize > 0 && f.ContentType != ""
}
