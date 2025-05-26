package entity

var Entities = map[string]any{
	(&RawApartment{}).TableName():       &RawApartment{},
	(&ProcessedApartment{}).TableName(): &ProcessedApartment{},
}
