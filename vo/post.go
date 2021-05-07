package vo

type CreatePostRequest struct{
	CategoryId uint `json:"category_id" binding:"required"`
	Title string `json:"title" biding:"required,max=10"`
	HeadImg string `json:"head_img"`
	Content string `json:"content" binding:"required"`
}