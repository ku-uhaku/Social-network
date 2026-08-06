"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";
import * as notificationsApi from "@/lib/api/notifications";

const NotificationContext = createContext(null);

// match backend
const ws_new_notification = "new_notification";
const ws_notification_expired = "notification_expired";
const ws_notification_read = "notification_read";

export function NotificationProvider({ children }) {
  const { user } = useAuth();
  const { subscribe } = useWebSocket();
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);

  const loadNotifications = async () => {
    if (!user) return;
    try {
      const data = await notificationsApi.getNotifications();
      const list = data?.data?.notifications ?? data?.notifications ?? [];
      const count = data?.data?.unread_count ?? data?.unread_count ?? 0;
      setNotifications(list);
      setUnreadCount(count);
    } catch {
      // keep existing state on failure
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user) {
      loadNotifications();
    } else {
      setNotifications([]);
      setUnreadCount(0);
      setLoading(true);
    }
  }, [user]);

  useEffect(() => {
    if (!user) return;

    const unsubNew = subscribe(ws_new_notification, (payload) => {
      if (!payload) return;
      setNotifications((prev) => {
        // deduplicate in case of notifications with same id
        // TODO: this is stupid
        if (prev.some((n) => n.id === payload.id)) return prev;
        return [payload, ...prev];
      });
      setUnreadCount((prev) => prev + 1);
    });

    const unsubExpired = subscribe(ws_notification_expired, (payload) => {
      if (!payload?.notification_id) return;
      const id = payload.notification_id;
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_expired: 1 } : n))
      );
      setUnreadCount((prev) => Math.max(0, prev - 1));
    });

    const unsubRead = subscribe(ws_notification_read, (payload) => {
      if (payload?.all) {
        setNotifications((prev) => prev.map((n) => ({ ...n, is_read: 1 })));
        setUnreadCount(0);
      } else if (payload?.notification_id) {
        const id = payload.notification_id;
        setNotifications((prev) =>
          prev.map((n) => (n.id === id ? { ...n, is_read: 1 } : n))
        );
        setUnreadCount((prev) => Math.max(0, prev - 1));
      }
    });

    return () => {
      unsubNew();
      unsubExpired();
      unsubRead();
    };
  }, [user, subscribe]);

  const markRead = async (notificationId) => {
    try {
      await notificationsApi.markNotificationRead(notificationId);
      setNotifications((prev) =>
        prev.map((n) => (n.id === notificationId ? { ...n, is_read: 1 } : n))
      );
      setUnreadCount((prev) => Math.max(0, prev - 1));
    } catch {
      // ignore
    }
  };

  const markAllRead = async () => {
    try {
      await notificationsApi.markAllNotificationsRead();
      setNotifications((prev) => prev.map((n) => ({ ...n, is_read: 1 })));
      setUnreadCount(0);
    } catch {
      // ignore
    }
  };

  const setMarkRead = async (notificationId) => {
    try {
      await notificationsApi.expireNotification(notificationId);
      setNotifications((prev) =>
        prev.map((n) => (n.id === notificationId ? { ...n, is_expired: 1 } : n))
      );
      setUnreadCount((prev) => Math.max(0, prev - 1));
    } catch {
      // ignore
    }
  };

  return (
    <NotificationContext.Provider
      value={{ notifications, unreadCount, loading, markRead, markAllRead, setMarkRead, refresh: loadNotifications }}
    >
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotifications() {
  const ctx = useContext(NotificationContext);
  if (!ctx) throw new Error("useNotifications must be used within a NotificationProvider");
  return ctx;
}