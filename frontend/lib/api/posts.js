import { apiFetch } from "./fetcher";

export function getFeed() {
  return apiFetch("/api/v1/posts/feed", { method: "GET" });
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
