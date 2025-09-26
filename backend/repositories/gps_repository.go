package repositories

import (
	"backend/models"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type GPSRepository struct {
	orm orm.Ormer
}

func NewGPSRepository() *GPSRepository {
	return &GPSRepository{
		orm: orm.NewOrm(),
	}
}

func (r *GPSRepository) Create(gpsData *models.GPSData) error {
	_, err := r.orm.Insert(gpsData)
	return err
}

func (r *GPSRepository) FindRecent(limit int, since time.Time) ([]models.GPSData, error) {
	var gpsData []models.GPSData
	_, err := r.orm.QueryTable(new(models.GPSData)).
		Filter("timestamp__gte", since).
		OrderBy("-timestamp").
		Limit(limit).
		All(&gpsData)
	return gpsData, err
}

func (r *GPSRepository) FindByVehicle(vehicleId string, limit int) ([]models.GPSData, error) {
	var gpsData []models.GPSData
	_, err := r.orm.QueryTable(new(models.GPSData)).
		Filter("vehicle_id", vehicleId).
		OrderBy("-timestamp").
		Limit(limit).
		All(&gpsData)
	return gpsData, err
}

func (r *GPSRepository) FindByRoad(roadId uint, since time.Time) ([]models.GPSData, error) {
	var gpsData []models.GPSData
	_, err := r.orm.QueryTable(new(models.GPSData)).
		Filter("road_segment_id", roadId).
		Filter("timestamp__gte", since).
		OrderBy("-timestamp").
		All(&gpsData)
	return gpsData, err
}
