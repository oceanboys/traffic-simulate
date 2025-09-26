package repositories

import (
	"backend/models"
	"github.com/beego/beego/v2/client/orm"
)

type AlertRepository struct {
	orm orm.Ormer
}

func NewAlertRepository() *AlertRepository {
	return &AlertRepository{
		orm: orm.NewOrm(),
	}
}

func (r *AlertRepository) Create(alert *models.TrafficAlert) error {
	_, err := r.orm.Insert(alert)
	return err
}

func (r *AlertRepository) GetByID(id uint) (*models.TrafficAlert, error) {
	alert := &models.TrafficAlert{ID: id}
	err := r.orm.Read(alert)
	return alert, err
}

func (r *AlertRepository) GetActive() ([]models.TrafficAlert, error) {
	var alerts []models.TrafficAlert
	_, err := r.orm.QueryTable(new(models.TrafficAlert)).
		Filter("resolved", false).
		OrderBy("-timestamp").
		All(&alerts)
	return alerts, err
}

func (r *AlertRepository) GetByRoad(roadID uint) ([]models.TrafficAlert, error) {
	var alerts []models.TrafficAlert
	_, err := r.orm.QueryTable(new(models.TrafficAlert)).
		Filter("road_segment_id", roadID).
		Filter("resolved", false).
		OrderBy("-timestamp").
		All(&alerts)
	return alerts, err
}

func (r *AlertRepository) Resolve(id uint) error {
	alert := &models.TrafficAlert{ID: id}
	alert.Resolved = true
	_, err := r.orm.Update(alert, "resolved")
	return err
}

func (r *AlertRepository) GetBySeverity(severity string) ([]models.TrafficAlert, error) {
	var alerts []models.TrafficAlert
	_, err := r.orm.QueryTable(new(models.TrafficAlert)).
		Filter("severity", severity).
		Filter("resolved", false).
		OrderBy("-timestamp").
		All(&alerts)
	return alerts, err
}
