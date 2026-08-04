import { apiFetch } from "./fetcher";

export function getFeed({ limit, cursor } = {}) {
  const params = new URLSearchParams();
  if (limit) params.set("limit", limit);
  if (cursor) params.set("cursor", cursor);
  const qs = params.toString();
  return apiFetch(`/api/v1/posts/feed${qs ? `?${qs}` : ""}`, { method: "GET" });
}

export function getPost(postId) {
  return apiFetch(`/api/v1/posts?post_id=${encodeURIComponent(postId)}`, { method: "GET" });
}

export function createPost(payload) {
  return apiFetch("/api/v1/posts", {
    method: "POST",
    body: payload,
  });
}

export function getComments(postId) {
  return apiFetch(`/api/v1/posts/comments?post_id=${encodeURIComponent(postId)}`, { method: "GET" });
}

export function createComment(payload) {
  return apiFetch("/api/v1/posts/comments", {
    method: "POST",
    body: payload,
  });
}
