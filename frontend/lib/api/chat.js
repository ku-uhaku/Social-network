import { apiFetch } from "./fetcher";

export function getConversations() {
  return apiFetch("/api/v1/chat/conversations", { method: "GET" });
}

export function getDirectHistory(userId) {
  return apiFetch(`/api/v1/chat/direct?user_id=${userId}`, { method: "GET" });
}
