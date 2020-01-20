package srt

// SearchParams is used to set search parameter
type SearchParams struct {
	Dep            string
	Arr            string
	Date           string
	Time           string
	IncludeSoldOut bool
}

// ReserveParams is used to set search parameter
type ReserveParams struct {
	Train           *Train
	Passengers      []*Passenger
	SpecialSeatOnly bool
}
