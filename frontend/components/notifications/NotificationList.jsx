"use client";

import { useNotifications } from "@/contexts/NotificationContext";
import NotificationItem from "@/components/notifications/NotificationItem";

export default function NotificationList({ open }) {
  const { notifications, loading } = useNotifications();

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