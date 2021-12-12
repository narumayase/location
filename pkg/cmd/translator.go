package cmd

import (
	"location/api"
	"location/pkg/db/model"
)

func toModel(b *api.BranchOffice) *model.BranchOffice {
	return &model.BranchOffice{
		ID:        b.ID,
		Longitude: b.Longitude,
		Latitude:  b.Latitude,
		Address:   b.Address,
	}
}

func toApi(b *model.BranchOffice) *api.BranchOffice {
	return &api.BranchOffice{
		ID:        b.ID,
		Longitude: b.Longitude,
		Latitude:  b.Latitude,
		Address:   b.Address,
	}
}

func toApis(bos []*model.BranchOffice) []*api.BranchOffice {
	var apis []*api.BranchOffice

	for _, bo := range bos {
		apis = append(apis, toApi(bo))
	}
	return apis
}
