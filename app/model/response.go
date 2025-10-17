package model

type MetaInfo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
	Pages  int    `json:"pages"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}

type PekerjaanResponse struct {
	Data []PekerjaanAlumni `json:"data"`
	Meta MetaInfo          `json:"meta"`
}
