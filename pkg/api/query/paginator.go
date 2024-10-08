package query

import "math"

const (
	MinPage        = 1   // MinPage represents the minimum allowed value for the pagination query's Page parameter.
	MinPerPage     = 1   // MinPerPage represents the minimum allowed value for the pagination query's PerPage parameter.
	DefaultPerPage = 10  // DefaultPerPage represents the default value for the pagination query's PerPage parameter.
	MaxPerPage     = 100 // MaxPerPage represents the maximum allowed value for the pagination query's PerPage parameter.
)

// Paginator represents the paginator parameters in a query.
type Paginator struct {
	// Page represents the current page number.
	Page int `query:"page"`

	// PerPage represents the number of items per page.
	PerPage int `query:"per_page"`
}

// NewPaginator creates and returns a new Paginator instance with MinPage and DefaultPerPage.
func NewPaginator() *Paginator {
	return &Paginator{
		Page:    MinPage,
		PerPage: DefaultPerPage,
	}
}

// Normalize ensures valid values for Page and PerPage in the pagination query.
// If query.PerPage is less than zero, it is set to `DefaultPerPage`.
// If query.Page is less than one, it is set to `MinPage`.
// The maximum allowed value for query.PerPage is `MaxPerPage`.
func (p *Paginator) Normalize() {
	p.Page = int(math.Max(float64(MinPage), float64(p.Page)))

	if p.PerPage == 0 {
		p.PerPage = DefaultPerPage
	} else {
		p.PerPage = int(math.Max(math.Min(float64(p.PerPage), float64(MaxPerPage)), float64(MinPerPage)))
	}
}
