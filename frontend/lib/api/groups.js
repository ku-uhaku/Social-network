import { apiFetch } from "./fetcher";

export function getAllGroups() {
  return apiFetch("/api/v1/groups/all", { method: "GET" });
}

export function getGroup(groupId) {
  return apiFetch(`/api/v1/groups/detail?id=${encodeURIComponent(groupId)}`, { method: "GET" });
}

export function createGroup(payload) {
  return apiFetch("/api/v1/groups", { method: "POST", body: payload });
}

export function getGroupFeed(groupId, { limit, cursor } = {}) {
  const params = new URLSearchParams();
  params.set("id", groupId);
  if (limit) params.set("limit", limit);
  if (cursor) params.set("cursor", cursor);
  return apiFetch(`/api/v1/groups/feed?${params.toString()}`, { method: "GET" });
}

export function joinGroup(groupId) {
  return apiFetch("/api/v1/groups/join", { method: "POST", body: { group_id: groupId } });
}

export function leaveGroup(groupId) {
  return apiFetch("/api/v1/groups/leave", { method: "POST", body: { group_id: groupId } });
}

export function inviteUsers(groupId, targetUserIds) {
  return apiFetch("/api/v1/groups/invite", {
    method: "POST",
    body: { group_id: groupId, target_user_ids: targetUserIds },
  });
}

export function getAllUsers() {
  return apiFetch("/api/v1/user/all", { method: "GET" });
}
