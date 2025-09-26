package models

import "time"

// TrafficAlert 交通告警模型
type TrafficAlert struct {
	ID          uint         `orm:"pk;auto"`
	AlertType   string       `orm:"size(20);index"`
	VehicleID   string       `orm:"size(20);null"`
	RoadSegment *RoadSegment `orm:"rel(fk)"`
	AlertValue  float64      `orm:"digits(10);decimals(2);null"`
	Message     string       `orm:"size(200)"`
	Severity    string       `orm:"size(10);default(medium);index"`
	Resolved    bool         `orm:"default(false);index"`
	Timestamp   time.Time    `orm:"type(datetime);index"`
	CreatedAt   time.Time    `orm:"auto_now_add;type(datetime)"`
}

func (t *TrafficAlert) TableName() string {
	return "traffic_alerts"
}

// IsHighSeverity 判断是否为高严重程度告警
func (t *TrafficAlert) IsHighSeverity() bool {
	return t.Severity == "high" || t.Severity == "critical"
}

// IsRecent 判断是否为最近告警（1小时内）
func (t *TrafficAlert) IsRecent() bool {
	return time.Since(t.Timestamp) <= time.Hour
}

// GetSeverityLevel 获取严重程度等级
func (t *TrafficAlert) GetSeverityLevel() int {
	switch t.Severity {
	case "low":
		return 1
	case "medium":
		return 2
	case "high":
		return 3
	case "critical":
		return 4
	default:
		return 2
	}
}

// Resolve 解决告警
func (t *TrafficAlert) Resolve() {
	t.Resolved = true
}
