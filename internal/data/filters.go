package data

import (
	"greenlight.aartchik.net/internal/validator"
	"strings"
)

type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafelist []string
}


type Metadata struct {
	Current_page int  `json:"current_page"`
	Page_size int	  `json:"page_size"`
	First_page int    `json:"first_page"`
	Last_page int     `json:"last_page"`
	Total_records int `json:"total_records"`
	
}

func (f Filters) limit() int {
	return f.PageSize
}
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *validator.Validator, filter *Filters) {


	v.Check(filter.Page >=1 && filter.Page <=1000000, "page", "must be in the range from 1 to 1000000")
	v.Check(filter.PageSize >= 1 && filter.PageSize <= 100, "page_size", "must be in the range from 1 to 100")
	v.Check(validator.PermittedValue(filter.Sort, filter.SortSafelist...), "sort", "invalid sort value")

}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}