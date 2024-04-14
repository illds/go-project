package postgresql

import (
	"GOHW-1/internal/db"
	"GOHW-1/internal/model"
	"context"
	"database/sql"
	"errors"
)

type PickUpPointRepo struct {
	db db.Database
}

func NewPickUpPoints(database db.Database) *PickUpPointRepo {
	return &PickUpPointRepo{db: database}
}

// Create creates an instance of pick-up point by given model.PickUpPoint
func (r *PickUpPointRepo) Create(ctx context.Context, pickUpPoint *model.PickUpPoint) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO pick_up_points(name, address, contact) VALUES ($1, $2, $3) RETURNING id;`,
		pickUpPoint.Name, pickUpPoint.Address, pickUpPoint.Contact).Scan(&id)
	return id, err
}

// GetByID return model.PickUpPoint by given ID
func (r *PickUpPointRepo) GetByID(ctx context.Context, id int64) (*model.PickUpPoint, error) {
	var pickUpPoint model.PickUpPoint
	if err := r.db.Get(ctx, &pickUpPoint, "SELECT id,name,address,contact FROM pick_up_points WHERE id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrObjectNotFound
		}
		return nil, err
	}
	return &pickUpPoint, nil
}

// List retrieves all pick-up points
func (r *PickUpPointRepo) List(ctx context.Context) ([]model.PickUpPoint, error) {
	var pickUpPoints []model.PickUpPoint
	if err := r.db.Select(ctx, &pickUpPoints, "SELECT id, name, address, contact FROM pick_up_points"); err != nil {
		return nil, err
	}
	return pickUpPoints, nil
}

// Update modifies details of an existing pick-up point
func (r *PickUpPointRepo) Update(ctx context.Context, id int64, updateData model.PickUpPoint) error {
	result, err := r.db.Exec(ctx, "UPDATE pick_up_points SET name=$2, address=$3, contact=$4 WHERE id=$1",
		id, updateData.Name, updateData.Address, updateData.Contact)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrObjectNotFound
	}
	return nil
}

// Delete removes a pick-up point by its ID.
func (r *PickUpPointRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.Exec(ctx, "DELETE FROM pick_up_points WHERE id=$1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrObjectNotFound
	}
	return nil
}
