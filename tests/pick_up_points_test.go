package tests

import (
	"GOHW-1/internal/controller"
	"GOHW-1/internal/model"
	"GOHW-1/internal/repository/postgresql"
	"GOHW-1/tests/fixtures"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

//
// Integration tests for handlers from HW-3
//

const (
	nonExistingIdStr = "-1"
)

func TestPickUpPointController_CreatePickUpPoint(t *testing.T) {
	repo := postgresql.NewPickUpPoints(tdb.DB)

	pickUpPointTest1 := fixtures.PickUpPoint().Valid().V()
	pickUpPointTest2 := fixtures.PickUpPoint().V()

	type fields struct {
		repo controller.PickUpPointsRepo
	}

	type args struct {
		ctx     context.Context
		request model.PickUpPoint
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "smoke test",
			fields:  fields{repo: repo},
			args:    args{ctx: context.Background(), request: pickUpPointTest1},
			want:    200,
			wantErr: assert.NoError,
		},
		{
			name:    "zero-value pick-up point test",
			fields:  fields{repo: repo},
			args:    args{ctx: context.Background(), request: pickUpPointTest2},
			want:    200,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &controller.PickUpPointController{
				Repo: tt.fields.repo,
			}
			_, got, err := controller.CreatePickUpPoint(tt.args.ctx, tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("create(%v, %v)", tt.args.ctx, tt.args.request)) {
				return
			}
			assert.Equalf(t, tt.want, got, "create(%v, %v)", tt.args.ctx, tt.args.request)
		})
	}
}

func TestPickUpPointController_GetJSONByID(t *testing.T) {
	repo := postgresql.NewPickUpPoints(tdb.DB)

	pickUpPointTest1 := fixtures.PickUpPoint().Valid().V()
	idTest1, _ := repo.Create(context.Background(), &pickUpPointTest1)
	pickUpPointTest1.ID = idTest1
	wantTest1, _ := json.Marshal(pickUpPointTest1)

	pickUpPointTest2 := fixtures.PickUpPoint().Valid().V()
	idTest2, _ := repo.Create(context.Background(), &pickUpPointTest2)
	pickUpPointTest2.ID = idTest2
	wantTest2, _ := json.Marshal(pickUpPointTest2)

	type fields struct {
		Repo controller.PickUpPointsRepo
	}
	type args struct {
		ctx   context.Context
		idStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "smoke test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), strconv.FormatInt(idTest1, 10)},
			want:    wantTest1,
			want1:   http.StatusOK,
			wantErr: assert.NoError,
		},
		{
			name:    "smoke test 2",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), strconv.FormatInt(idTest2, 10)},
			want:    wantTest2,
			want1:   http.StatusOK,
			wantErr: assert.NoError,
		},
		{
			name:    "non existing ID test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), nonExistingIdStr},
			want:    nil,
			want1:   http.StatusNotFound,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &controller.PickUpPointController{
				Repo: tt.fields.Repo,
			}
			got, got1, err := controller.GetJSONByID(tt.args.ctx, tt.args.idStr)
			if !tt.wantErr(t, err, fmt.Sprintf("GetJSONByID(%v, %v)", tt.args.ctx, tt.args.idStr)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetJSONByID(%v, %v)", tt.args.ctx, tt.args.idStr)
			assert.Equalf(t, tt.want1, got1, "GetJSONByID(%v, %v)", tt.args.ctx, tt.args.idStr)
		})
	}
}

func TestPickUpPointController_UpdateByID(t *testing.T) {
	repo := postgresql.NewPickUpPoints(tdb.DB)

	pickUpPointTest1 := fixtures.PickUpPoint().Valid().V()
	idTest1, _ := repo.Create(context.Background(), &pickUpPointTest1)
	pickUpPointUpdatedTest1 := fixtures.PickUpPoint().Valid().ID(idTest1).Name("Maksim").V()
	wantTest1, _ := json.Marshal(pickUpPointUpdatedTest1)

	pickUpPointUpdatedTest2 := fixtures.PickUpPoint().Valid().V()

	type fields struct {
		Repo controller.PickUpPointsRepo
	}
	type args struct {
		ctx   context.Context
		idStr string
		unm   model.PickUpPoint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "smoke test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), strconv.FormatInt(idTest1, 10), pickUpPointUpdatedTest1},
			want:    wantTest1,
			want1:   http.StatusOK,
			wantErr: assert.NoError,
		},
		{
			name:    "non existing ID test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), nonExistingIdStr, pickUpPointUpdatedTest2},
			want:    nil,
			want1:   http.StatusNotFound,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &controller.PickUpPointController{
				Repo: tt.fields.Repo,
			}
			got, got1, err := controller.UpdateByID(tt.args.ctx, tt.args.idStr, tt.args.unm)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateByID(%v, %v, %v)", tt.args.ctx, tt.args.idStr, tt.args.unm)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateByID(%v, %v, %v)", tt.args.ctx, tt.args.idStr, tt.args.unm)
			assert.Equalf(t, tt.want1, got1, "UpdateByID(%v, %v, %v)", tt.args.ctx, tt.args.idStr, tt.args.unm)
		})
	}
}

func TestPickUpPointController_DeleteByID(t *testing.T) {
	repo := postgresql.NewPickUpPoints(tdb.DB)

	pickUpPointTest1 := fixtures.PickUpPoint().Valid().V()
	idTest1, _ := repo.Create(context.Background(), &pickUpPointTest1)

	type fields struct {
		Repo controller.PickUpPointsRepo
	}
	type args struct {
		ctx   context.Context
		idStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "smoke test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), strconv.FormatInt(idTest1, 10)},
			want:    http.StatusOK,
			wantErr: assert.NoError,
		},
		{
			name:    "non existing ID test",
			fields:  fields{Repo: repo},
			args:    args{context.Background(), nonExistingIdStr},
			want:    http.StatusNotFound,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &controller.PickUpPointController{
				Repo: tt.fields.Repo,
			}
			got, err := controller.DeleteByID(tt.args.ctx, tt.args.idStr)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteByID(%v, %v)", tt.args.ctx, tt.args.idStr)) {
				return
			}
			assert.Equalf(t, tt.want, got, "DeleteByID(%v, %v)", tt.args.ctx, tt.args.idStr)
		})
	}
}
