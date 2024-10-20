package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

// type Coordinates struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// }

type GeoJSONData struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"` // This can be Coordinates or []Coordinates for different shape types
}

func (g GeoJSONData) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GeoJSONData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &g)
}

type GeoData struct {
	gorm.Model
	UserID uint           `gorm:"not null"` // Foreign key to User
	Data   datatypes.JSON `gorm:"type:jsonb;not null"`
}

func (GeoData) TableName() string {
	return "geo_data"
}
