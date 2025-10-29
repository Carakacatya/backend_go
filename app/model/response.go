package model

type MetaInfo struct {
	Page   int    `bson:"page" json:"page"`
	Limit  int    `bson:"limit" json:"limit"`
	Total  int    `bson:"total" json:"total"`
	Pages  int    `bson:"pages" json:"pages"`
	SortBy string `bson:"sort_by" json:"sortBy"`
	Order  string `bson:"order" json:"order"`
	Search string `bson:"search" json:"search"`
}
