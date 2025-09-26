package utils

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

// CreateTables 创建数据库表
func CreateTables() error {
	o := orm.NewOrm()

	// 创建路段表
	_, err := o.Raw(`
	CREATE TABLE IF NOT EXISTS road_segments (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	start_lng DECIMAL(10,6) NOT NULL,
	start_lat DECIMAL(10,6) NOT NULL,
	end_lng DECIMAL(10,6) NOT NULL,
	end_lat DECIMAL(10,6) NOT NULL,
	max_speed INT DEFAULT 60,
	capacity INT DEFAULT 1000,
	length DECIMAL(8,2),
	road_type VARCHAR(10) DEFAULT 'urban',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	INDEX idx_name (name),
	INDEX idx_road_type (road_type)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Exec()

	if err != nil {
		logs.Error("创建路段表失败: ", err)
		return err
	}

	// 创建GPS数据表
	_, err = o.Raw(`
	CREATE TABLE IF NOT EXISTS gps_data (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	vehicle_id VARCHAR(20) NOT NULL,
	longitude DECIMAL(10,6) NOT NULL,
	latitude DECIMAL(10,6) NOT NULL,
	speed INT DEFAULT 0,
	direction INT,
	timestamp DATETIME NOT NULL,
	road_segment_id INT UNSIGNED,
	vehicle_type VARCHAR(10) DEFAULT 'car',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	INDEX idx_vehicle_id (vehicle_id),
	INDEX idx_timestamp (timestamp),
	INDEX idx_road_segment (road_segment_id),
	FOREIGN KEY (road_segment_id) REFERENCES road_segments(id) ON DELETE SET NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Exec()

	if err != nil {
		logs.Error("创建GPS数据表失败: ", err)
		return err
	}

	// 创建交通告警表
	_, err = o.Raw(`
	CREATE TABLE IF NOT EXISTS traffic_alerts (
	id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	alert_type VARCHAR(50) NOT NULL,
	vehicle_id VARCHAR(20),
	road_segment_id INT UNSIGNED,
	alert_value DECIMAL(10,2),
	message TEXT,
	severity VARCHAR(20) DEFAULT 'medium',
	timestamp DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	INDEX idx_alert_type (alert_type),
	INDEX idx_vehicle_id (vehicle_id),
	INDEX idx_road_segment (road_segment_id),
	INDEX idx_timestamp (timestamp),
	FOREIGN KEY (road_segment_id) REFERENCES road_segments(id) ON DELETE SET NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`).Exec()

	if err != nil {
		logs.Error("创建交通告警表失败: ", err)
		return err
	}

	logs.Info("数据库表创建成功")
	return nil
}

// DropTables 删除数据库表
func DropTables() error {
	o := orm.NewOrm()

	// 删除表（注意外键约束顺序）
	tables := []string{"traffic_alerts", "gps_data", "road_segments"}

	for _, table := range tables {
		_, err := o.Raw("DROP TABLE IF EXISTS " + table).Exec()
		if err != nil {
			logs.Error("删除表 "+table+" 失败: ", err)
			return err
		}
	}

	logs.Info("数据库表删除成功")
	return nil
}

// ResetTables 重置数据库表
func ResetTables() error {
	// 先删除表
	if err := DropTables(); err != nil {
		return err
	}

	// 再创建表
	return CreateTables()
}
