package utils

import (
	"errors"
	"math"
)

type MetaData struct {
	CurrentPage  int
	perPage      int
	TotalPages   int
	TotalRecords int
}

func ComputeMetaData(CurrentPage, perPage, totalRecords int) (MetaData, int, error) {

	//calculating offset

	offset := (CurrentPage - 1) * perPage
	//calculating total pages

	totalPages := math.Ceil(float64(totalRecords) / float64(perPage))

	if totalPages == 0 || CurrentPage > int(totalPages) {
		return MetaData{}, -1, errors.New("Oops ! no records found")
	}
	return MetaData{
		CurrentPage:  CurrentPage,
		perPage:      perPage,
		TotalPages:   int(totalPages),
		TotalRecords: totalRecords,
	}, offset, nil
}
