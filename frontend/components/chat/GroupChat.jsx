"use client";

import { useEffect, useRef, useState } from "react";
import { useWebSocket } from "@/contexts/WebSocketContext";
import { useGroupChat } from "@/contexts/GroupChatContext";
import { getGroupHistory } from "@/lib/api/chat";
import { formatMessageTime, displayName } from "@/lib/utils";
import Avatar from "@/components/shared/Avatar";
import Composer from "./Composer";
import "@/css/chat.css";

const pageSize = 30;

export default function GroupChat({ groupId, title, meId, onClose }) {
  const { send, subscribe } = useWebSocket();
  const { openGroup, closeGroup } = useGroupChat();

  useEffect(() => {
    openGroup(groupId);
    return () => closeGroup(groupId);
  }, [groupId, openGroup, closeGroup]);

  const [messages, setMessages] = useState([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(true);
  const [loadingOlder, setLoadingOlder] = useState(false);
  const [draft, setDraft] = useState("");
  const listEndRef = useRef(null);

  useEffect(() => {
    let cancelled = false;
    getGroupHistory(groupId, { page: 1 })
      .then((res) => {
        if (cancelled) return;
        const history = Array.isArray(res?.data) ? res.data : [];
        setMessages(history.slice().reverse());
        setHasMore(history.length === pageSize);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [groupId]);

  useEffect(() => {
    const unsub = subscribe("new_group_message", (msg) => {
      if (!msg || msg.group_id !== groupId) return;
      setMessages((prev) => (prev.some((m) => m.id === msg.id) ? prev : [...prev, msg]));
    });
    return unsub;
  }, [subscribe, groupId, meId]);

  useEffect(() => {
    listEndRef.current?.scrollIntoView({ block: "end" });
  }, [messages.length]);

  const loadOlder = async () => {
    if (loadingOlder || !hasMore) return;
    setLoadingOlder(true);
    try {
      const nextPage = page + 1;
      const res = await getGroupHistory(groupId, { page: nextPage });
      const older = Array.isArray(res?.data) ? res.data : [];
      setMessages((prev) => [...older.slice().reverse(), ...prev]);
      setPage(nextPage);
      setHasMore(older.length === pageSize);
      setLoadingOlder(false);
    } catch {
      setLoadingOlder(false);
    }
  };

  const sendMessage = () => {
    const content = draft.trim();
    if (!content) return;
    send("send_group_message", { group_id: groupId, content });
    setDraft("");
  };

  return (
    <div className="chatModalOverlay" onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className="chatModal chatThread">
        <div className="chatPanelHeader">
          <div className="chatThreadTitle">
            <strong className="chatContactName">{title}</strong>
          </div>
          <button type="button" className="chatCloseButton" onClick={onClose}>
            &times;
          </button>
        </div>

        <div className="chatMessages">
          {loading && <p className="chatEmpty">Loading...</p>}
          {!loading && hasMore && (
            <button
              type="button"
              className="chatLoadOlder"
              onClick={loadOlder}
              disabled={loadingOlder}
            >
              {loadingOlder ? "Loading..." : "Load previous"}
            </button>
          )}
          {!loading && messages.length === 0 && (
            <p className="chatEmpty">No messages yet. Say hello!</p>
          )}
          {messages.map((m) => {
            const mine = m.sender_id === meId;
            return (
              <div key={m.id} className="chatGroupMessage">
                {!mine && (
                  <div className="chatSender">
                    <Avatar avatar={m.avatar} username={m.username} size={22} />
                    <span className="chatSenderName">{displayName(m)}</span>
                  </div>
                )}
                <div className={`chatBubble ${mine ? "mine" : "theirs"}`}>
                  <span className="chatBubbleText">{m.content}</span>
                  <span className="chatBubbleTime">{formatMessageTime(m.created_at)}</span>
                </div>
              </div>
            );
          })}
          <div ref={listEndRef} />
        </div>

        <Composer draft={draft} setDraft={setDraft} onSend={sendMessage} />
      </div>
    </div>
  );
}
