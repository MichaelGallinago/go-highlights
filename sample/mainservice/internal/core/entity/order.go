package entity

type OrderEvent struct {
	Id     int64  `json:"order_id"`
	UserId int64  `json:"user_id"`
	Status string `json:"status"`
}

type Order struct {
	UserId int64
	Items  []Item
}

type Item struct {
	ProductId int
	Quantity  int64
}

func (o *Order) GetOrderItems() []int {
	res := make([]int, len(o.Items))
	for i, item := range o.Items {
		res[i] = item.ProductId
	}

	return res
}
