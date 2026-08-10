import { apiFetch } from "./fetcher";

export function getNotifications({ limit = 20, lastId = 0 } = {}) {
  return apiFetch(`/api/v1/notifications?limit=${limit}&last_id=${lastId}`, { method: "GET" });
}

export function markNotificationRead(notificationId) {
  return apiFetch("/api/v1/notifications/read", {
    method: "POST",
    body: { notification_id: notificationId },
  });
}

export function markAllNotificationsRead() {
  return apiFetch("/api/v1/notifications/read", {
    method: "POST",
    body: { all: true },
  });
}

export function expireNotification(notificationId) {
  return apiFetch("/api/v1/notifications/expire", {
    method: "POST",
    body: { notification_id: notificationId },
  });
}