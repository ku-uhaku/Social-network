package handler

import (
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

// GetFeed GET /api/v1/posts/feed
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		helper.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	posts, err := h.Service.GetFeed(r.Context(), user.ID)
	if err != nil {
		helper.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(w, http.StatusOK, "Feed retrieved successfully", posts)
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
