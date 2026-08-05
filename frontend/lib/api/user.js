import { apiFetch } from "./fetcher";

export function getUserProfile(username) {
  return apiFetch(`/api/v1/user/profile?username=${encodeURIComponent(username)}`, { method: "GET" });
}

export function getUserPosts(username) {
  return apiFetch(`/api/v1/user/posts?username=${encodeURIComponent(username)}`, { method: "GET" });
}

export function updateProfilePrivacy(isPublic) {
  return apiFetch("/api/v1/user/profile/update", {
    method: "PUT",
    body: { is_public: isPublic ? 1 : 0 },
  });
}

export function followUser(targetUserId) {
  return apiFetch("/api/v1/user/follow", {
    method: "POST",
    body: { target_user_id: targetUserId },
  });
}

export function unfollowUser(targetUserId) {
  return apiFetch("/api/v1/user/unfollow", {
    method: "POST",
    body: { target_user_id: targetUserId },
  });
}

export function getFollowers(userId) {
  return apiFetch(`/api/v1/user/followers?id=${userId}`, { method: "GET" });
}

export function getFollowing(userId) {
  return apiFetch(`/api/v1/user/following?id=${userId}`, { method: "GET" });
}