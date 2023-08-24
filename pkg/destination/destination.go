package destination

import (
	"errors"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Destination interface {
	UploadFiles(fileList []string, prefix string) error
}

func GetDestinationByName(name string) (Destination, error) {
	var destination Destination

	var destinationFile types.DestinationFile
	var dest types.Destination
	file, err := os.Open(os.Getenv("DESTINATION_FILE_PATH"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &destinationFile)

	for _, d := range destinationFile.Destinations {
		if d.Name == name {
			dest = d
		}
	}

	switch dest.Type {
	case types.Aws:
		destination = NewS3BucketDestination(S3DestinationParams{
			Bucket: dest.Bucket,
			Region: dest.Region,
		})
	default:
		return nil, errors.New("destination Type don't match")
	}

	return destination, nil
}
