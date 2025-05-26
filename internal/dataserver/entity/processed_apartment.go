package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ProcessedApartments []ProcessedApartment

type ProcessedApartment struct {
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

func (*ProcessedApartment) TableName() string {
	return "processed_apartments"
}

// TODO: uncomment after implementing processor
func (a *ProcessedApartment) CreateFromPayload(data []byte) error {
	if err := a.LoadFromPayload(data); err != nil {
		return err
	}

	return a.Create()
}

func (a *ProcessedApartment) LoadFromPayload(data []byte) error {
	return json.Unmarshal(data, a)
}

func (a *ProcessedApartment) Create() error {
	return Db().Create(a).Error
}

func NewProcessedApartment() *ProcessedApartment {
	return &ProcessedApartment{}
}

func FindProcessedApartment(uuid string) *ProcessedApartment {
	if uuid == "" {
		return nil
	}

	result := &ProcessedApartment{}

	if err := Db().Where("uuid = ?", uuid).First(result).Error; err != nil {
		return nil
	}

	return result
}

func ProcessedApartmentsCount() int64 {
	var count int64
	Db().Model(&ProcessedApartment{}).Count(&count)
	return count
}
