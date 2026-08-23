package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"kuu/internal/helper"
	"kuu/internal/middleware"
	"kuu/internal/requests"
	"kuu/internal/service"
)

// CreatePost POST /api/v1/posts
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	payload, err := requests.ParseCreatePostPayload(r)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if errs := requests.ValidateCreatePost(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	post, err := h.Service.CreatePost(r.Context(), user.ID, payload)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	helper.Success(w, http.StatusCreated, "Post created successfully", post)
}

// GetFeed GET /api/v1/posts/feed?limit=10&cursor=123
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed < 1 {
			helper.Error(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsed
		if limit > 50 {
			limit = 50
		}
	}

	var cursor *int64
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		parsed, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || parsed < 1 {
			helper.Error(w, http.StatusBadRequest, "Invalid cursor parameter")
			return
		}
		cursor = &parsed
	}

	posts, hasMore, err := h.Service.GetFeed(r.Context(), user.ID, limit, cursor)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var nextCursor *int64
	if hasMore && len(posts) > 0 {
		lastID := posts[len(posts)-1].ID
		nextCursor = &lastID
	}

	helper.Success(w, http.StatusOK, "Feed retrieved successfully", map[string]interface{}{
		"posts":       posts,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}

// GetPost GET /api/v1/posts?post_id=123
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID, err := helper.GetParamInt64(r, "post_id")
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := h.Service.GetPost(r.Context(), user.ID, postID)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			helper.Error(w, http.StatusNotFound, "Post not found")
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.Success(w, http.StatusOK, "Post retrieved successfully", post)
}

// CreateComment POST /api/v1/posts/comments
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	payload, err := requests.ParseCreateCommentPayload(r)
	if err != nil {
		helper.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// file, header, err := r.FormFile("image")
	// if err != nil {
	// 	http.Error(w, "failed to get uploaded file", http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()

	// _, err = helper.SaveUploadedImage(file, header)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// fmt.Fprintln(w, path)
	if errs := requests.ValidateCreateComment(payload); len(errs) > 0 {
		helper.WriteJSON(w, http.StatusUnprocessableEntity, false, "Validation failed", nil, errs)
		return
	}

	comment, err := h.Service.AddComment(r.Context(), user.ID, payload)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusCreated, "Comment added successfully", comment)
}

// GetComments GET /api/v1/posts/comments?post_id=123
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil || postID <= 0 {
		helper.Error(w, http.StatusBadRequest, "Invalid post_id parameter")
		return
	}

	comments, err := h.Service.GetComments(r.Context(), user.ID, postID)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			helper.Error(w, http.StatusForbidden, err.Error())
			return
		}
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Comments retrieved successfully", comments)
}
