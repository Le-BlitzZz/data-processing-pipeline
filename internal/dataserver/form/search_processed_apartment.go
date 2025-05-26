package form

type SearchProcessedApartment struct {
	Count int    `form:"count"`
	Order string `form:"order"`
}
