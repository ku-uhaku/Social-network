"use client";

import Avatar from "@/components/shared/Avatar";
import { useNotifications } from "@/contexts/NotificationContext";
import { acceptFollowRequest, declineFollowRequest } from "@/lib/api/user";

function timeAgo(dateString) {
  const date = new Date(dateString);
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  if (seconds < 60) return "now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

export default function NotificationItem({ notification }) {
  const { markRead, expireNotification } = useNotifications();
  const actions = notification.actions || {};

  const handleAction = async (action) => {
    const api = action === "accept" ? acceptFollowRequest : declineFollowRequest;
    try {
      await api(notification.actor_id);
      await expireNotification(notification.id); // expire
    } catch {
      // ignore
    }
  };

  return (
    <div className={`notificationItem ${notification.is_read ? "" : "unread"} ${notification.is_expired ? "expired" : ""}`}>
      <div className="notificationItemHeader">
        <Avatar
          avatar={notification.actor_avatar}
          username={notification.actor_username || "user"}
          size={36}
        />
        <div className="notificationMeta">
          <strong className="notificationTitle">{notification.title}</strong>
          <span className="notificationTime">{timeAgo(notification.created_at)}</span>
        </div>
        {!notification.is_read && !notification.is_expired && (
          <button
            type="button"
            className="notificationReadButton"
            title="Mark as read"
            onClick={() => markRead(notification.id)}
          >
            ✓
          </button>
        )}
      </div>
      <p className="notificationMessage">{notification.message}</p>
      {!notification.is_expired && actions.buttons && (
        <div className="notificationActions">
          {actions.buttons.map((btn) => (
            <button
              key={btn.action}
              type="button"
              className={`notificationAction ${btn.action}`}
              onClick={() => handleAction(btn.action)}
            >
              {btn.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
