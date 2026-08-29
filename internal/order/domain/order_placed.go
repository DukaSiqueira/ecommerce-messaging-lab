package domain

import "time"

type OrderPlaced struct {
	OrderID    string
	EventID    string
	Items      []OrderPlacedItem
	CustomerID string
	OccurredAt time.Time
}

type OrderPlacedItem struct {
	ProductID string
	Quantity  int
}
