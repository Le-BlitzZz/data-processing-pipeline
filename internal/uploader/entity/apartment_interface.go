package entity

type Apartment interface {
	CreateFromPayload(data []byte) error
}

type ApartmentFactory func() Apartment

func NewRawApartment() Apartment {
	return &RawApartment{}
}

func NewProcessedApartment() Apartment {
	return &ProcessedApartment{}
}
