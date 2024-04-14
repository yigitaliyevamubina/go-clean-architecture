package v1

import (
	"fourth-exam/post-service-clean-arch/internal/entity"
	"fourth-exam/post-service-clean-arch/internal/usecase"
	"fourth-exam/post-service-clean-arch/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type postRoutes struct {
	t usecase.Post
	l logger.Interface
}

func newPostRoutes(handler *gin.RouterGroup, t usecase.Post, l logger.Interface) {
	r := &postRoutes{t, l}

	h := handler.Group("/post")
	{
		h.POST("/create", r.CreatePost)
		h.PUT("/update/:id", r.UpdatePost)
		h.GET("/:id", r.GetPostById)
		h.PUT("/like", r.LikePost)
		h.PUT("/dislike", r.DislikePost)
		h.DELETE("/delete/:id", r.DeletePost)
	}

	handler.GET("/posts/:page/:limit/:user_id", r.ListPostsByUserId)
	handler.GET("/posts/:page/:limit", r.ListPosts)
}

// CreatePost
// @Router /post/create [post]
// @Summary create post
// @Tags Post
// @Description Insert a new post with provided details
// @Accept json
// @Produce json
// @Param PostDetails body entity.Post true "Create post"
// @Success 201 {object} entity.Post
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) CreatePost(c *gin.Context) {
	var (
		body       entity.Post
		jspMarshal protojson.MarshalOptions
	)
	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if err != nil {
		p.l.Error(err, "http - v1 - create post")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	body.Id = uuid.New().String()

	post, err := p.t.CreatePost(c.Request.Context(), &body)
	if err != nil {
		p.l.Error(err, "http - v1 - create post")
		errorResponse(c, http.StatusInternalServerError, "create post service problems")

		return
	}

	c.JSON(http.StatusOK, post)
}

// Update Post
// @Router /post/update/{id} [put]
// @Summary update post
// @Tags Post
// @Description Update post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param PostInfo body entity.Post true "Update Post"
// @Success 201 {object} entity.Post
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) UpdatePost(c *gin.Context) {
	var (
		body        entity.Post
		jspbMarshal protojson.MarshalOptions
	)
	id := c.Param("id")

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		p.l.Error(err, "http - v1 - update post")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	body.Id = id
	response, err := p.t.UpdatePost(c.Request.Context(), &body)
	if err != nil {
		p.l.Error(err, "http - v1 - update post")
		errorResponse(c, http.StatusInternalServerError, "update post service problems")

		return
	}

	c.JSON(http.StatusOK, response)
}

// Like Post
// @Router /post/like [put]
// @Summary like post
// @Tags Post
// @Description Like post
// @Accept json
// @Produce json
// @Param post_id body entity.PostRequest true "Like Post"
// @Success 201 {object} entity.Post
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) LikePost(c *gin.Context) {
	var (
		body        entity.PostRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		p.l.Error(err, "http - v1 - like post")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	post, err := p.t.GetPost(c.Request.Context(), body.PostId)
	if err != nil {
		p.l.Error(err, "http - v1 - get post")
		errorResponse(c, http.StatusInternalServerError, "get post service problems")

		return
	}
	post.Likes += 1

	response, err := p.t.UpdatePost(c.Request.Context(), post)
	if err != nil {
		p.l.Error(err, "http - v1 - update post")
		errorResponse(c, http.StatusInternalServerError, "update post service problems")

		return
	}

	c.JSON(http.StatusOK, response)
}

// DisLike Post
// @Router /post/dislike [put]
// @Summary dislike post
// @Tags Post
// @Description Dislike post
// @Accept json
// @Produce json
// @Param post_id body entity.PostRequest true "Dislike Post"
// @Success 201 {object} entity.Post
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) DislikePost(c *gin.Context) {
	var (
		body        entity.PostRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		p.l.Error(err, "http - v1 - dislike post")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	post, err := p.t.GetPost(c.Request.Context(), body.PostId)
	if err != nil {
		p.l.Error(err, "http - v1 - get post")
		errorResponse(c, http.StatusInternalServerError, "get post service problems")

		return
	}
	post.Likes -= 1

	response, err := p.t.UpdatePost(c.Request.Context(), post)
	if err != nil {
		p.l.Error(err, "http - v1 - update post")
		errorResponse(c, http.StatusInternalServerError, "update post service problems")

		return
	}

	c.JSON(http.StatusOK, response)
}

// Get Post By Id
// @Router /post/{id} [get]
// @Summary get post by id
// @Tags Post
// @Description Get post
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} entity.Post
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) GetPostById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	post, err := p.t.GetPost(c.Request.Context(), id)
	if err != nil {
		p.l.Error(err, "http - v1 - get post")
		errorResponse(c, http.StatusInternalServerError, "get post service problems")

		return
	}

	c.JSON(http.StatusOK, post)
}

// Delete Post
// @Router /post/delete/{id} [delete]
// @Summary delete post
// @Tags Post
// @Description Delete post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} entity.MessageResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) DeletePost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	err := p.t.DeletePost(c.Request.Context(), id)
	if err != nil {
		p.l.Error(err, "http - v1 - delete post")
		errorResponse(c, http.StatusInternalServerError, "delete post service problems")

		return
	}

	c.JSON(http.StatusOK, entity.MessageResponse{
		Message: "post was successfully deleted",
	})
}

// Get All Posts
// @Router /posts/{page}/{limit} [get]
// @Summary get all posts
// @Tags Post
// @Description get all posts
// @Accept json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Param orderBy query string false "orderBy" Enums(content, title, category) "Order by"
// @Success 201 {object} entity.Posts
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) ListPosts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	var (
		req entity.GetListFilter
	)
	orderBy := c.Query("orderBy")
	if orderBy != "" {
		req.OrderBy = orderBy
	} else {
		req.OrderBy = "created_at"
	}

	page := c.Param("page")
	pageToInt, err := strconv.Atoi(page)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts - parsing page to int")
		errorResponse(c, http.StatusInternalServerError, "list posts page parse error")

		return
	}

	limit := c.Param("limit")
	LimitToInt, err := strconv.Atoi(limit)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts - parsing limit to int")
		errorResponse(c, http.StatusInternalServerError, "list posts limit parse error")

		return
	}

	req.Page = int64(pageToInt)
	req.Limit = int64(LimitToInt)

	posts, err := p.t.ListPosts(c.Request.Context(), &req)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts")
		errorResponse(c, http.StatusInternalServerError, " list post service problems")

		return
	}

	c.JSON(http.StatusOK, posts)
}

// Get All Posts by user id
// @Router /posts/{page}/{limit}/{user_id} [get]
// @Summary get all posts
// @Tags Post
// @Description get all posts by user id
// @Accept json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Param orderBy query string false "orderBy" Enums(content, title, category, created_at, updated_at) "Order by"
// @Param user_id path string true "user_id"
// @Success 201 {object} entity.Posts
// @Failure 400 {object} response
// @Failure 500 {object} response
func (p *postRoutes) ListPostsByUserId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	var (
		req entity.GetListFilter
	)
	orderBy := c.Query("orderBy")
	if orderBy != "" {
		req.OrderBy = orderBy
	} else {
		req.OrderBy = "created_at"
	}

	userId := c.Param("user_id")
	req.UserId = userId

	page := c.Param("page")
	pageToInt, err := strconv.Atoi(page)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts - parsing page to int")
		errorResponse(c, http.StatusInternalServerError, "list posts page parse error")

		return
	}

	limit := c.Param("limit")
	LimitToInt, err := strconv.Atoi(limit)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts - parsing limit to int")
		errorResponse(c, http.StatusInternalServerError, "list posts limit parse error")

		return
	}

	req.Page = int64(pageToInt)
	req.Limit = int64(LimitToInt)

	posts, err := p.t.ListPosts(c.Request.Context(), &req)
	if err != nil {
		p.l.Error(err, "http - v1 - list posts by user id")
		errorResponse(c, http.StatusInternalServerError, " list post by user id service problems")

		return
	}

	c.JSON(http.StatusOK, posts)
}
