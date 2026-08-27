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
      const hasChoice = n.actions?.buttons?.some(
        (b) => b.action === "accept" || b.action === "decline"
      );
      if (!n.is_read && !hasChoice) markRead(n.id);
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