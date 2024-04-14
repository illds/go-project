package storage

import (
	"GOHW-1/internal/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func setUp(t *testing.T, data []model.PickUpPoint) string {
	t.Helper()
	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "pickUpPoints*.json")
	if err != nil {
		t.Fatalf("Unable to create temporary file: %v", err)
	}
	filename := tmpFile.Name()

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Unable to marshal pick up points data: %v", err)
	}
	if _, err := tmpFile.Write(jsonData); err != nil {
		t.Fatalf("Unable to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Unable to close temporary file: %v", err)
	}

	return filename
}

func readFile(pickUpPointsFileName string) ([]model.PickUpPoint, error) {
	file, err := os.Open(pickUpPointsFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open the file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get the information about the file")
	}

	// Is file empty
	if fileInfo.Size() == 0 {
		return []model.PickUpPoint{}, nil
	}

	rawBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var pickUpPoints []model.PickUpPoint
	if err := json.Unmarshal(rawBytes, &pickUpPoints); err != nil {
		return nil, err
	}

	return pickUpPoints, nil
}

func tearDown(filename string) {
	os.Remove(filename)
}
