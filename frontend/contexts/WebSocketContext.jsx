"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { getWebSocketUrl } from "@/lib/sockets";

const WebSocketContext = createContext(null);

const reconnect_delay = 1000;

export function WebSocketProvider({ children }) {
  const { user } = useAuth();
  const [connected, setConnected] = useState(false);
  const wsRef = useRef(null);
  const reconnectTimerRef = useRef(null);
  const listenersRef = useRef(new Map()); // type -> Set<callback>
  const userRef = useRef(user);

  useEffect(() => {
    userRef.current = user;
  }, [user]);

  const connect = () => {
    if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
    if (!userRef.current || wsRef.current) return;

    const ws = new WebSocket(getWebSocketUrl());
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
    };

    ws.onmessage = (event) => {
      let data;
      try {
        data = JSON.parse(event.data);
      } catch {
        return;
      }
      const callbacks = listenersRef.current.get(data.type);
      callbacks?.forEach((cb) => cb(data.payload, data));
    };

    ws.onclose = () => {
      setConnected(false);
      wsRef.current = null;
      scheduleReconnect();
    };

    ws.onerror = () => {
      ws.close();
    };
  };

  const scheduleReconnect = () => {
    if (!userRef.current) return;
    reconnectTimerRef.current = setTimeout(connect, reconnect_delay);
  };

  useEffect(() => {
    if (user) {
      connect();
    }
    return () => {
      if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
      if (wsRef.current) wsRef.current.close();
      wsRef.current = null;
    };
  }, [user]);

  const send = (type, payload) => {
    const ws = wsRef.current;
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(JSON.stringify({ type, payload }));
  };

  const subscribe = (type, callback) => {
    if (!listenersRef.current.has(type)) {
      listenersRef.current.set(type, new Set());
    }
    listenersRef.current.get(type).add(callback);
    return () => {
      const set = listenersRef.current.get(type);
      if (set) {
        set.delete(callback);
        if (set.size === 0) listenersRef.current.delete(type);
      }
    };
  };

  return (
    <WebSocketContext.Provider value={{ connected, send, subscribe }}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  const ctx = useContext(WebSocketContext);
  if (!ctx) throw new Error("useWebSocket must be used within a WebSocketProvider");
  return ctx;
}
