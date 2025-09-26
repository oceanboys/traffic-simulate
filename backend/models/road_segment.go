package models

import "time"

// RoadSegment 路段模型
type RoadSegment struct {
	ID        uint      `orm:"pk;auto"`
	Name      string    `orm:"size(100);index"`
	StartLng  float64   `orm:"digits(10);decimals(6)"`
	StartLat  float64   `orm:"digits(10);decimals(6)"`
	EndLng    float64   `orm:"digits(10);decimals(6)"`
	EndLat    float64   `orm:"digits(10);decimals(6)"`
	MaxSpeed  int       `orm:"default(60)"`
	Capacity  int       `orm:"default(1000)"`
	Length    float64   `orm:"digits(8);decimals(2);null"`
	RoadType  string    `orm:"size(10);default(urban);index"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (r *RoadSegment) TableName() string {
	return "road_segments"
}

// GetCenterPoint 获取路段中心点
func (r *RoadSegment) GetCenterPoint() (float64, float64) {
	return (r.StartLng + r.EndLng) / 2, (r.StartLat + r.EndLat) / 2
}

// GetLength 计算路段长度（简化计算）
func (r *RoadSegment) GetLength() float64 {
	if r.Length > 0 {
		return r.Length
	}
	// 简单的距离计算（实际应该使用更精确的地理计算）
	latDiff := r.EndLat - r.StartLat
	lngDiff := r.EndLng - r.StartLng
	return (latDiff*latDiff + lngDiff*lngDiff) * 111.0 // 粗略转换为公里
}

// IsVehicleInSegment 判断车辆是否在路段内
func (r *RoadSegment) IsVehicleInSegment(lng, lat float64) bool {
	// 简化的判断逻辑，实际应该使用更精确的地理算法
	return lng >= r.StartLng && lng <= r.EndLng && lat >= r.StartLat && lat <= r.EndLat
}
