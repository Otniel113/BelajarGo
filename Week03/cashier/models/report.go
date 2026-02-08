package models

type Report struct {
	TotalRevenue      int             `json:"total_revenue"`
	TotalTransaction  int             `json:"total_transaction"`
	MostSoldProduct   MostSoldProduct `json:"most_sold_product"`
}

type MostSoldProduct struct {
	Name    string `json:"name"`
	SoldQty int    `json:"sold_qty"`
}
