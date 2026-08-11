"use client";

import { useEffect } from "react";
import { useNotifications } from "@/contexts/NotificationContext";
import NotificationItem from "@/components/notifications/NotificationItem";

export default function NotificationList({ open }) {
  const { notifications, markRead, loading } = useNotifications();

  // consume on open
  useEffect(() => {
    if (!open) return;
    notifications.forEach((n) => {
      if (n.type === "group_event_created" && !n.is_read) {
        markRead(n.id);
      }
    });
  }, [open, notifications, markRead]);

  return (
    <section className={`notificationPanel ${open ? "open" : ""}`}>
      {loading ? (
        <p className="notificationEmpty">Loading...</p>
      ) : notifications.length === 0 ? (
        <p className="notificationEmpty">Empty</p>
      ) : (
        <div className="notificationList">
          {notifications.map((n) => (
            <NotificationItem key={n.id} notification={n} />
          ))}
        </div>
      )}
    </section>
  );
}