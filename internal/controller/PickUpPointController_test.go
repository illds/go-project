package controller

import (
	"GOHW-1/internal/model"
	"GOHW-1/tests/fixtures"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.Background()
		id          = int64(1)
		idStr       = "1"
		nonValidStr = "non-valid"
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().GetByID(gomock.Any(), id).Return(fixtures.PickUpPoint().Valid().P(), nil)

		// act
		result, status, _ := s.pickUpPointController.GetJSONByID(ctx, idStr)

		// assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"Name\":\"Ildus\",\"Address\":\"Saint-P\",\"Contact\":\"123\"}", string(result))
	})
	t.Run("NoRows table test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().GetByID(gomock.Any(), id).Return(nil, sql.ErrNoRows)

		// act
		result, status, err := s.pickUpPointController.GetJSONByID(ctx, idStr)

		// assert
		require.Equal(t, http.StatusNotFound, status)
		require.Equal(t, "pick-up point not found", err.Error())
		assert.Equal(t, "", string(result))
	})
	t.Run("ErrObjectNotFound test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().GetByID(gomock.Any(), id).Return(nil, model.ErrObjectNotFound)

		// act
		result, status, err := s.pickUpPointController.GetJSONByID(ctx, idStr)

		// assert
		require.Equal(t, http.StatusNotFound, status)
		require.Equal(t, "error occured: not found", err.Error())
		assert.Equal(t, "", string(result))
	})
	t.Run("Non-valid idStr test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		// act
		result, status, err := s.pickUpPointController.GetJSONByID(ctx, nonValidStr)

		// assert
		require.Equal(t, http.StatusBadRequest, status)
		require.Equal(t, "id must be a number", err.Error())
		assert.Equal(t, "", string(result))
	})
}

func Test_Create(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(0)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Create(gomock.Any(), fixtures.PickUpPoint().Valid().ID(id).P()).Return(id, nil)

		// act
		result, status, err := s.pickUpPointController.CreatePickUpPoint(ctx, fixtures.PickUpPoint().Valid().V())

		// assert
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, nil, err)
		assert.Equal(t, "{\"ID\":0,\"Name\":\"Ildus\",\"Address\":\"Saint-P\",\"Contact\":\"123\"}", string(result))
	})
	t.Run("fail test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Create(gomock.Any(), fixtures.PickUpPoint().Valid().ID(0).P()).Return(id, fmt.Errorf("some error"))

		// act
		result, status, err := s.pickUpPointController.CreatePickUpPoint(ctx, fixtures.PickUpPoint().Valid().V())

		// assert
		require.Equal(t, http.StatusInternalServerError, status)
		require.Equal(t, "can not created pick-up point", err.Error())
		assert.Equal(t, "", string(result))
	})
}

func Test_Update(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.Background()
		id          = int64(1)
		zeroValue   = int64(0)
		idStr       = "1"
		nonValidStr = "non-valid"
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Update(gomock.Any(), id, fixtures.PickUpPoint().Valid().ID(zeroValue).V()).Return(nil)

		// act
		result, status, err := s.pickUpPointController.UpdateByID(ctx, idStr, fixtures.PickUpPoint().Valid().ID(zeroValue).V())

		// assert
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, nil, err)
		assert.Equal(t, "{\"ID\":1,\"Name\":\"Ildus\",\"Address\":\"Saint-P\",\"Contact\":\"123\"}", string(result))
	})
	t.Run("fail test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Update(gomock.Any(), id, fixtures.PickUpPoint().Valid().ID(0).V()).Return(fmt.Errorf("some error"))

		// act
		result, status, err := s.pickUpPointController.UpdateByID(ctx, idStr, fixtures.PickUpPoint().Valid().ID(zeroValue).V())

		// assert
		require.Equal(t, http.StatusNotFound, status)
		require.Equal(t, "pick-up point not found", err.Error())
		assert.Equal(t, "", string(result))
	})
	t.Run("Non-valid idStr test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		// act
		result, status, err := s.pickUpPointController.UpdateByID(ctx, nonValidStr, fixtures.PickUpPoint().Valid().V())

		// assert
		require.Equal(t, http.StatusBadRequest, status)
		require.Equal(t, "id must be a number", err.Error())
		assert.Equal(t, "", string(result))
	})
}

func Test_Delete(t *testing.T) {
	t.Parallel()
	var (
		ctx         = context.Background()
		id          = int64(1)
		idStr       = "1"
		nonValidStr = "non-valid"
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Delete(gomock.Any(), id).Return(nil)

		// act
		status, err := s.pickUpPointController.DeleteByID(ctx, idStr)

		// assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, nil, err)
	})
	t.Run("fail test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockPickUpPoints.EXPECT().Delete(gomock.Any(), id).Return(fmt.Errorf("some error"))

		// act
		status, err := s.pickUpPointController.DeleteByID(ctx, idStr)

		// assert
		require.Equal(t, http.StatusNotFound, status)
		assert.Equal(t, "some error", err.Error())
	})
	t.Run("Non-valid idStr test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()

		// act
		status, err := s.pickUpPointController.DeleteByID(ctx, nonValidStr)

		// assert
		require.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, "id must be a number", err.Error())
	})
}
