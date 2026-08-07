"use client";

import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";
import * as notificationsApi from "@/lib/api/notifications";

const NotificationContext = createContext(null);

export function NotificationProvider({ children }) {
  const { user } = useAuth();
  const { subscribe } = useWebSocket();
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [prevUserId, setPrevUserId] = useState(user?.id);

  // Reset state when the logged-in user changes (covers logout).
  if (user?.id !== prevUserId) {
    setPrevUserId(user?.id);
    setNotifications([]);
    setUnreadCount(0);
    setLoading(true);
  }

  const markOne = useCallback((id, patch) => {
    setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, ...patch } : n)));
    setUnreadCount((prev) => Math.max(0, prev - 1));
  }, []);

  const markAllReadLocal = useCallback(() => {
    setNotifications((prev) => prev.map((n) => ({ ...n, is_read: 1 })));
    setUnreadCount(0);
  }, []);

  useEffect(() => {
    if (!user) return;

    let cancelled = false;
    notificationsApi.getNotifications()
      .then((data) => {
        if (cancelled) return;
        setNotifications(data?.data?.notifications ?? data?.notifications ?? []);
        setUnreadCount(data?.data?.unread_count ?? data?.unread_count ?? 0);
      })
      .catch(() => {
        // keep existing state on failure
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [user]);


  // match backend
  useEffect(() => {
    if (!user) return;

    const unsubNew = subscribe("new_notification", (payload) => {
      if (!payload) return;
      setNotifications((prev) => {
        // deduplicate in case of notifications with same id
        // TODO: this is stupid
        if (prev.some((n) => n.id === payload.id)) return prev;
        return [payload, ...prev];
      });
      setUnreadCount((prev) => prev + 1);
    });

    const unsubExpired = subscribe("notification_expired", (payload) => {
      if (payload?.notification_id) markOne(payload.notification_id, { is_expired: 1 });
    });

    const unsubRead = subscribe("notification_read", (payload) => {
      if (payload?.all) {
        markAllReadLocal();
      } else if (payload?.notification_id) {
        markOne(payload.notification_id, { is_read: 1 });
      }
    });

    return () => {
      unsubNew();
      unsubExpired();
      unsubRead();
    };
  }, [user, subscribe, markOne, markAllReadLocal]);

  const markRead = async (notificationId) => {
    try {
      await notificationsApi.markNotificationRead(notificationId);
      markOne(notificationId, { is_read: 1 });
    } catch {
      // ignore
    }
  };

  const markAllRead = async () => {
    try {
      await notificationsApi.markAllNotificationsRead();
      markAllReadLocal();
    } catch {
      // ignore
    }
  };

  const expireNotification = async (notificationId) => {
    try {
      await notificationsApi.expireNotification(notificationId);
      markOne(notificationId, { is_expired: 1 });
    } catch {
      // ignore
    }
  };

  return (
    <NotificationContext.Provider
      value={{ notifications, unreadCount, loading, markRead, markAllRead, expireNotification }}
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
