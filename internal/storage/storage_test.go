package storage

import (
	"GOHW-1/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage_PickUpPointsRead(t *testing.T) {
	t.Parallel()
	t.Run("smoke test", func(t *testing.T) {
		// arrange
		sampleData := []model.PickUpPoint{
			{ID: 1, Name: "Point A", Address: "123 A St", Contact: "111-222-3333"},
			{ID: 2, Name: "Point B", Address: "456 B St", Contact: "444-555-6666"},
		}
		filename := setUp(t, sampleData)
		defer tearDown(filename)
		storage := &Storage{}
		pickUpPointsFileName = filename

		// act
		pickUpPoints, err := storage.PickUpPointsRead()

		// assert
		require.Equal(t, err, nil)
		assert.Equal(t, sampleData, pickUpPoints)
	})
	t.Run("can not open file test", func(t *testing.T) {
		// arrange
		var sampleData []model.PickUpPoint
		filename := setUp(t, sampleData)
		tearDown(filename)
		storage := &Storage{}
		pickUpPointsFileName = filename

		// act
		pickUpPoints, err := storage.PickUpPointsRead()

		// assert
		require.Equal(t, "unable to open the file", err.Error())
		assert.Equal(t, sampleData, pickUpPoints)
	})
}

func TestStorage_WritePickUpPoints(t *testing.T) {
	t.Parallel()
	t.Run("smoke test", func(t *testing.T) {
		// arrange
		pickUpPoint := model.PickUpPoint{ID: 1, Name: "Point A", Address: "123 A St", Contact: "111-222-3333"}
		sampleData := []model.PickUpPoint{}
		filename := setUp(t, sampleData)
		defer tearDown(filename)
		storage := &Storage{}
		pickUpPointsFileName = filename

		// act
		err := storage.writePickUpPoints([]model.PickUpPoint{pickUpPoint})

		// assert
		pickUpPointsResult, _ := readFile(filename)
		require.Equal(t, nil, err)
		assert.Equal(t, pickUpPointsResult, []model.PickUpPoint{pickUpPoint})
	})
	t.Run("smoke test 2", func(t *testing.T) {
		// arrange
		pickUpPoint := model.PickUpPoint{ID: 1, Name: "Point A", Address: "123 A St", Contact: "111-222-3333"}
		pickUpPoint2 := model.PickUpPoint{ID: 2, Name: "Point B", Address: "456 B St", Contact: "444-555-6666"}
		sampleData := []model.PickUpPoint{}
		filename := setUp(t, sampleData)
		defer tearDown(filename)
		storage := &Storage{}
		pickUpPointsFileName = filename

		// act
		err := storage.writePickUpPoints([]model.PickUpPoint{pickUpPoint, pickUpPoint2})

		// assert
		pickUpPointsResult, _ := readFile(filename)
		require.Equal(t, nil, err)
		assert.Equal(t, []model.PickUpPoint{pickUpPoint, pickUpPoint2}, pickUpPointsResult)
	})
}
