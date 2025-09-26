package models

import "time"

// GPSData GPS数据模型
type GPSData struct {
	ID          uint         `orm:"pk;auto"`
	VehicleID   string       `orm:"size(20);index"`
	Longitude   float64      `orm:"digits(10);decimals(6)"`
	Latitude    float64      `orm:"digits(10);decimals(6)"`
	Speed       int          `orm:"default(0)"`
	Direction   int          `orm:"null"`
	Timestamp   time.Time    `orm:"type(datetime);index"`
	RoadSegment *RoadSegment `orm:"rel(fk);null"`
	VehicleType string       `orm:"size(10);default(car)"`
	CreatedAt   time.Time    `orm:"auto_now_add;type(datetime)"`
}

func (g *GPSData) TableName() string {
	return "gps_data"
}

// GetLocation 获取位置信息
func (g *GPSData) GetLocation() (float64, float64) {
	return g.Longitude, g.Latitude
}

// IsSpeeding 判断是否超速
func (g *GPSData) IsSpeeding(maxSpeed int) bool {
	return g.Speed > maxSpeed
}
