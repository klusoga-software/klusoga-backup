package destination

type Destination interface {
	UploadFiles(fileList []string, prefix string) error
}
