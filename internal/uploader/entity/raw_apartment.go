package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type RawApartment struct {
	gorm.Model

	UUID                 string  `json:"id"`
	City                 string  `json:"city"`
	Type                 string  `json:"type"`
	SquareMeters         float64 `json:"squareMeters,string"`
	Rooms                float64 `json:"rooms,string"`
	Floor                float64 `json:"floor,string"`
	FloorCount           float64 `json:"floorCount,string"`
	BuildYear            float64 `json:"buildYear,string"`
	Latitude             float64 `json:"latitude,string"`
	Longitude            float64 `json:"longitude,string"`
	CentreDistance       float64 `json:"centreDistance,string"`
	PoiCount             float64 `json:"poiCount,string"`
	SchoolDistance       float64 `json:"schoolDistance,string"`
	ClinicDistance       float64 `json:"clinicDistance,string"`
	PostOfficeDistance   float64 `json:"postOfficeDistance,string"`
	KindergartenDistance float64 `json:"kindergartenDistance,string"`
	RestaurantDistance   float64 `json:"restaurantDistance,string"`
	CollegeDistance      float64 `json:"collegeDistance,string"`
	PharmacyDistance     float64 `json:"pharmacyDistance,string"`
	Ownership            string  `json:"ownership"`
	HasParkingSpace      string  `json:"hasParkingSpace"`
	HasBalcony           string  `json:"hasBalcony"`
	HasElevator          string  `json:"hasElevator"`
	HasSecurity          string  `json:"hasSecurity"`
	HasStorageRoom       string  `json:"hasStorageRoom"`
	Price                int64   `json:"price,string"`

	Split string `gorm:"type:enum('train','test','val')" json:"split"`
}

func (*RawApartment) TableName() string {
	return "raw_apartments"
}

func (a *RawApartment) CreateFromPayload(data []byte) error {
	if err := a.LoadFromPayload(data); err != nil {
		return err
	}

	return a.Create()
}

func (a *RawApartment) LoadFromPayload(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *RawApartment) Create() error {
	return Db().Create(a).Error
}
