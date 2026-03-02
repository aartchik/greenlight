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