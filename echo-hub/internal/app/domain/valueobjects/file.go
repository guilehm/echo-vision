package valueobjects

type File struct {
	filepath    string
	filename    string
	filesize    int64
	contentType string
}

func NewFile(filepath, filename, contentType string, filesize int64) *File {
	return &File{
		filepath:    filepath,
		filename:    filename,
		filesize:    filesize,
		contentType: contentType,
	}
}

func (f *File) Filepath() string {
	return f.filepath
}

func (f *File) Filename() string {
	return f.filename
}

func (f *File) Filesize() int64 {
	return f.filesize
}

func (f *File) ContentType() string {
	return f.contentType
}

func (f *File) IsValid() bool {
	return f.filepath != "" && f.filename != "" && f.filesize > 0 && f.contentType != ""
}
