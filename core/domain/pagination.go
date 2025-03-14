package domain

type Pagination struct {
	Search    string `json:"search" form:"search"`
	Page      uint   `json:"page" form:"page"`
	Limit     uint   `json:"limit" form:"limit"`
	Offset    uint   `json:"-"`
	Sort      string `json:"sort" form:"sort"`
	TotalRow  uint   `json:"total_rows,omitempty"`
	TotalPage uint   `json:"total_page,omitempty"`
}
