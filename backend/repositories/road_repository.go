package repositories

import (
	"backend/models"
	"github.com/beego/beego/v2/client/orm"
)

type RoadRepository struct {
	orm orm.Ormer
}

func NewRoadRepository() *RoadRepository {
	return &RoadRepository{
		orm: orm.NewOrm(),
	}
}

func (r *RoadRepository) GetAll() ([]models.RoadSegment, error) {
	var segments []models.RoadSegment
	_, err := r.orm.QueryTable(new(models.RoadSegment)).All(&segments)
	return segments, err
}

func (r *RoadRepository) GetByID(id uint) (*models.RoadSegment, error) {
	segment := &models.RoadSegment{ID: id}
	err := r.orm.Read(segment)
	return segment, err
}

func (r *RoadRepository) Create(segment *models.RoadSegment) error {
	_, err := r.orm.Insert(segment)
	return err
}

func (r *RoadRepository) Update(segment *models.RoadSegment) error {
	_, err := r.orm.Update(segment)
	return err
}

func (r *RoadRepository) Delete(id uint) error {
	_, err := r.orm.Delete(&models.RoadSegment{ID: id})
	return err
}
