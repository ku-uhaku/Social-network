"use client";

import { createContext, useCallback, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";

const GroupChatContext = createContext(null);

export function GroupChatProvider({ children }) {
  const { user } = useAuth();
  const { subscribe } = useWebSocket();
  const [unread, setUnread] = useState({});
  const activeGroupRef = useRef(null);
  const prevUserIdRef = useRef(user?.id);

  useEffect(() => {
    const prev = prevUserIdRef.current;
    prevUserIdRef.current = user?.id;
    if (user?.id !== prev) {
      setUnread({});
      activeGroupRef.current = null;
    }
  }, [user?.id]);

  useEffect(() => {
    if (!user) return;
    const unsub = subscribe("new_group_message", (msg) => {
      if (!msg || msg.sender_id === user.id) return;
      if (activeGroupRef.current === msg.group_id) return;
      setUnread((prev) => ({ ...prev, [msg.group_id]: (prev[msg.group_id] || 0) + 1 }));
    });
    return unsub;
  }, [user, subscribe]);

  const openGroup = useCallback((groupId) => {
    activeGroupRef.current = groupId;
    setUnread((prev) => (prev[groupId] ? { ...prev, [groupId]: 0 } : prev));
  }, []);

  const closeGroup = useCallback((groupId) => {
    if (activeGroupRef.current === groupId) activeGroupRef.current = null;
  }, []);

  return (
    <GroupChatContext.Provider value={{ unread, openGroup, closeGroup }}>
      {children}
    </GroupChatContext.Provider>
  );
}

export function useGroupChat() {
  const ctx = useContext(GroupChatContext);
  if (!ctx) throw new Error("useGroupChat must be used within a GroupChatProvider");
  return ctx;
}
