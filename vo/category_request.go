package vo

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`    // name字段是required
}