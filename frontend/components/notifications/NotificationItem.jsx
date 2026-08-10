"use client";

import { useRouter } from "next/navigation";
import Avatar from "@/components/shared/Avatar";
import { useNotifications } from "@/contexts/NotificationContext";
import { acceptFollowRequest, declineFollowRequest } from "@/lib/api/user";
import {
  acceptGroupInvitation,
  declineGroupInvitation,
  acceptJoinRequest,
  declineJoinRequest,
} from "@/lib/api/groups";

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
  const router = useRouter();
  const { markRead, expireNotification } = useNotifications();
  const actions = notification.actions || {};

  const handleAction = async (action) => {
    const typeActions = {
      follow_request: {
        accept: (n) => acceptFollowRequest(n.actor_id),
        decline: (n) => declineFollowRequest(n.actor_id),
      },
      group_invitation: {
        accept: (n) => acceptGroupInvitation(n.payload?.group_id),
        decline: (n) => declineGroupInvitation(n.payload?.group_id),
      },
      group_join_request: {
        accept: (n) => acceptJoinRequest(n.payload?.group_id, n.actor_id),
        decline: (n) => declineJoinRequest(n.payload?.group_id, n.actor_id),
      },
      group_event_created: {
        view: (n) => router.push(`/group/${n.payload?.group_id}/events`),
      },
    };
    const handler = typeActions[notification.type]?.[action];
    if (!handler) return;
    try {
      await handler(notification);
      if (notification.type === "group_event_created") {
        await markRead(notification.id);
      } else {
        await expireNotification(notification.id); // expire
      }
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
