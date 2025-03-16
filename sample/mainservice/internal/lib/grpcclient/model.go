package grpcclient

import "csu-lessons/product/api"

type ProductState struct {
	ProductId int64
	Count     int64
}

func fromProtoModel(states []*api.ProductState) []ProductState {
	result := make([]ProductState, len(states))
	for i, state := range states {
		result[i] = ProductState{
			ProductId: state.ProductId,
			Count:     state.Count,
		}
	}

	return result
}
