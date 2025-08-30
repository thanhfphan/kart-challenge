package dto

type IDRequest struct {
	ID int64 `form:"id" uri:"id" json:"id" binding:"required"`
}
