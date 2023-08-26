package destination

import (
	"errors"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/types"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
)

type Destination interface {
	UploadFiles(fileList []string) error
	String() string
}

func getDestinationFileData() (*types.DestinationFile, error) {
	var destinationFile types.DestinationFile
	file, err := getDestinationFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(data, &destinationFile)

	return &destinationFile, nil
}

func getDestinationFile() (*os.File, error) {
	destinationFilePath := os.Getenv("DESTINATION_FILE_PATH")
	if destinationFilePath == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		destinationFilePath = path.Join(homedir, ".klusoga-backup", "destinations.yaml")
		checkDestinationFilePath(destinationFilePath)
	}

	file, err := os.OpenFile(destinationFilePath, os.O_RDWR|os.O_CREATE, 0750)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func checkDestinationFilePath(path string) {
	os.MkdirAll(filepath.Dir(path), 0750)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, _ := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0750)
		destinationFile := types.DestinationFile{}
		data, _ := yaml.Marshal(destinationFile)
		file.Write(data)
		file.Close()
	}
}

func GetDestinationByName(name string) (Destination, error) {
	var destination Destination
	var dest types.Destination

	destinationFile, err := getDestinationFileData()
	if err != nil {
		return nil, err
	}

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
			Name:   dest.Name,
		})
	default:
		return nil, errors.New("destination Type don't match")
	}

	return destination, nil
}

func ListDestinations() ([]Destination, error) {
	var destinations []Destination

	destinationFile, err := getDestinationFileData()
	if err != nil {
		return nil, err
	}

	for _, d := range destinationFile.Destinations {
		switch d.Type {
		case types.Aws:
			destinations = append(destinations, NewS3BucketDestination(S3DestinationParams{
				Bucket: d.Bucket,
				Region: d.Region,
				Name:   d.Name,
			}))
		}
	}

	return destinations, nil
}

func AddDestination(destination types.Destination) error {
	destinationFile, err := getDestinationFileData()
	if err != nil {
		return err
	}

	destinationFile.Destinations = append(destinationFile.Destinations, destination)

	file, err := getDestinationFile()
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(destinationFile)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()

	return nil
}

func DeleteDestination(name string) error {
	destinationFile, err := getDestinationFileData()
	if err != nil {
		return err
	}

	destinationFile.Destinations = slices.DeleteFunc(destinationFile.Destinations, func(destination types.Destination) bool {
		return destination.Name == name
	})

	file, err := getDestinationFile()
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := yaml.Marshal(destinationFile)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()
	return nil
}
