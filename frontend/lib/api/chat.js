import { apiFetch } from "./fetcher";

export function getConversations() {
  return apiFetch("/api/v1/chat/conversations", { method: "GET" });
}

export function getDirectHistory(userId, { beforeId = 0, limit = 30 } = {}) {
  const params = new URLSearchParams({ user_id: String(userId), limit: String(limit) });
  if (beforeId > 0) params.set("before_id", String(beforeId));
  return apiFetch(`/api/v1/chat/direct?${params}`, { method: "GET" });
}

export function markChatRead(userId) {
  return apiFetch("/api/v1/chat/read", {
    method: "POST",
    body: { user_id: userId },
  });
}

export function getGroupHistory(groupId, { page = 1 } = {}) {
  const params = new URLSearchParams({ group_id: String(groupId), page: String(page) });
  return apiFetch(`/api/v1/chat/group?${params}`, { method: "GET" });
}
