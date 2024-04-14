package entity

type Post struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	Likes     int64  `json:"likes"`
	Dislikes  int64  `json:"dislikes"`
	Views     int64  `json:"views"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetListFilter struct {
	Page    int64  `json:"page"`
	Limit   int64  `json:"limit"`
	OrderBy string `json:"order_by"`
	UserId  string `json:"user_id"`
}

type Posts struct {
	Count int64   `json:"count"`
	Items []*Post `json:"posts"`
}

type PostRequest struct {
	PostId string `json:"post_id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
