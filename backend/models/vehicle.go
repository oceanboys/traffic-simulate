package models

import "time"

// Vehicle 车辆模型
type Vehicle struct {
	ID          uint      `json:"id" orm:"auto;pk"`
	VehicleID   string    `json:"vehicle_id" orm:"size(50);unique"`
	X           float64   `json:"x" orm:"column(x_coordinate)"`
	Y           float64   `json:"y" orm:"column(y_coordinate)"`
	Speed       float64   `json:"speed"`
	Direction   float64   `json:"direction"`
	VehicleType string    `json:"vehicle_type" orm:"size(20)"`
	Status      string    `json:"status" orm:"size(20)"`
	CreatedAt   time.Time `json:"created_at" orm:"auto_now_add"`
	UpdatedAt   time.Time `json:"updated_at" orm:"auto_now"`
}

// TableName 返回表名
func (v *Vehicle) TableName() string {
	return "vehicles"
}
