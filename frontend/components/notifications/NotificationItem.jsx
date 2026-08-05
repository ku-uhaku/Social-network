"use client";

import { resolveMediaSrc } from "@/lib/utils";
import { useNotifications } from "@/contexts/NotificationContext";

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
  const { markRead, setMarkRead } = useNotifications();
  const avatarSrc = resolveMediaSrc(notification.actor_avatar);
  const actions = notification.actions || {};

  const handleAction = (action) => {
    if (action === "accept" || action === "decline") {
      setMarkRead(notification.id); // expire
    }
  };

  return (
    <div className={`notificationItem ${notification.is_read ? "" : "unread"} ${notification.is_expired ? "expired" : ""}`}>
      <div className="notificationItemHeader">
        {avatarSrc ? (
          <img className="notificationAvatar" src={avatarSrc} alt={notification.actor_username || "user"} />
        ) : (
          <span className="notificationAvatarFallback">{notification.actor_username || "?"}</span>
        )}
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