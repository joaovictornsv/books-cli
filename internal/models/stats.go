package models

type LibraryStats struct {
	Year             int
	ByStatus         map[string]int
	ByCategory       map[string]int
	FinishedThisYear int
	PriorityWishlist int
}
