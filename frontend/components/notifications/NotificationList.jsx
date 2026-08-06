"use client";

import { useNotifications } from "@/contexts/NotificationContext";
import NotificationItem from "@/components/notifications/NotificationItem";

export default function NotificationList({ open }) {
  const { notifications, unreadCount, markAllRead, loading } = useNotifications();

  return (
    <section className={`notificationPanel ${open ? "open" : ""}`}>
      <div className="notificationPanelHeader">
        <h3 className="notificationPanelTitle">
          <img
            className="notificationPanelBell"
            src={unreadCount > 0 ? "/images/notification_on.png" : "/images/notification_off.png"}
            alt=""
          />
          Notifications
        </h3>
        {unreadCount > 0 && (
          <button type="button" className="notificationMarkAllRead" onClick={markAllRead}>
            Mark all read
          </button>
        )}
      </div>

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