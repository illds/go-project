package controller

import (
	controller "GOHW-1/internal/controller/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

type pickUpPointRepoFixtures struct {
	ctrl                  *gomock.Controller
	pickUpPointController PickUpPointController
	mockPickUpPoints      *controller.MockPickUpPointsRepo
}

func setUp(t *testing.T) pickUpPointRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPickUpPoints := controller.NewMockPickUpPointsRepo(ctrl)
	pickUpPointController := PickUpPointController{Repo: mockPickUpPoints}
	return pickUpPointRepoFixtures{
		ctrl:                  ctrl,
		pickUpPointController: pickUpPointController,
		mockPickUpPoints:      mockPickUpPoints,
	}
}

func (a *pickUpPointRepoFixtures) tearDown() {
	a.ctrl.Finish()
}
