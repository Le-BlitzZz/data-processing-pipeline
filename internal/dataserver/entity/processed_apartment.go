package entity

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ProcessedApartments []ProcessedApartment

type ProcessedApartment struct {
	gorm.Model

	UUID string `json:"id"`

	SquareMeters   float64 `json:"num__squareMeters,string"`
	Rooms          float64 `json:"num__rooms,string"`
	Floor          float64 `json:"num__floor,string"`
	FloorCount     float64 `json:"num__floorCount,string"`
	CentreDistance float64 `json:"num__centreDistance,string"`
	PoiCount       float64 `json:"num__poiCount,string"`
	Age            float64 `json:"num__age,string"`

	HasParkingSpace float64 `json:"bool__hasParkingSpace,string"`
	HasBalcony      float64 `json:"bool__hasBalcony,string"`
	HasElevator     float64 `json:"bool__hasElevator,string"`
	HasSecurity     float64 `json:"bool__hasSecurity,string"`
	HasStorageRoom  float64 `json:"bool__hasStorageRoom,string"`
	HasSchool       float64 `json:"bool__hasSchool,string"`
	HasClinic       float64 `json:"bool__hasClinic,string"`
	HasPostoffice   float64 `json:"bool__hasPostoffice,string"`
	HasKindergarten float64 `json:"bool__hasKindergarten,string"`
	HasRestaurant   float64 `json:"bool__hasRestaurant,string"`
	HasCollege      float64 `json:"bool__hasCollege,string"`
	HasPharmacy     float64 `json:"bool__hasPharmacy,string"`

	CityBialystok   float64 `json:"cat__city_bialystok,string"`
	CityBydgoszcz   float64 `json:"cat__city_bydgoszcz,string"`
	CityCzestochowa float64 `json:"cat__city_czestochowa,string"`
	CityGdansk      float64 `json:"cat__city_gdansk,string"`
	CityGdynia      float64 `json:"cat__city_gdynia,string"`
	CityKatowice    float64 `json:"cat__city_katowice,string"`
	CityKrakow      float64 `json:"cat__city_krakow,string"`
	CityLodz        float64 `json:"cat__city_lodz,string"`
	CityLublin      float64 `json:"cat__city_lublin,string"`
	CityPoznan      float64 `json:"cat__city_poznan,string"`
	CityRadom       float64 `json:"cat__city_radom,string"`
	CityRzeszow     float64 `json:"cat__city_rzeszow,string"`
	CitySzczecin    float64 `json:"cat__city_szczecin,string"`
	CityWarszawa    float64 `json:"cat__city_warszawa,string"`
	CityWroclaw     float64 `json:"cat__city_wroclaw,string"`

	TypeApartmentBuilding float64 `json:"cat__type_apartmentBuilding,string"`
	TypeBlockOfFlats      float64 `json:"cat__type_blockOfFlats,string"`
	TypeTenement          float64 `json:"cat__type_tenement,string"`

	Price int64 `json:"price,string"`

	Split string `gorm:"type:enum('train','test','val')" json:"split"`
}

func (*ProcessedApartment) TableName() string {
	return "processed_apartments"
}

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
