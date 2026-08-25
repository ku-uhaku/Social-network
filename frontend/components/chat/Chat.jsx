"use client";

import { useEffect, useRef, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";
import { useAudio } from "@/contexts/AudioContext";
import { getConversations, getDirectHistory, markChatRead } from "@/lib/api/chat";
import { formatMessageTime, displayName } from "@/lib/utils";
import Avatar from "@/components/shared/Avatar";
import Composer from "./Composer";
import "@/css/chat.css";

const pageSize = 30;

export function useChat() {
  const { user } = useAuth();
  const { send, subscribe } = useWebSocket();
  const { playSfx } = useAudio();
  const playSfxRef = useRef(playSfx);
  useEffect(() => {
    playSfxRef.current = playSfx;
  }, [playSfx]);

  const [open, setOpenState] = useState(false);
  const [conversations, setConversations] = useState([]);
  const [unread, setUnread] = useState({}); // { [userId]: count }

  // Only the open conversation is kept in memory; history is refetched on open.
  const [activeId, setActiveId] = useState(null);
  const [messages, setMessages] = useState([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const [loading, setLoading] = useState(false);
  const [loadingOlder, setLoadingOlder] = useState(false);
  const [draft, setDraft] = useState("");

  // Conversation list, with the unread counts the server remembers.
  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    getConversations()
      .then((res) => {
        if (cancelled) return;
        const list = Array.isArray(res?.data) ? res.data : [];
        setConversations(list);
        setUnread(Object.fromEntries(list.map((c) => [c.user_id, c.unread_count || 0])));
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [user]);

  // Realtime DMs: append to the open thread, otherwise bump its badge.
  useEffect(() => {
    if (!user) return;
    const unsub = subscribe("new_direct_message", (msg) => {
      if (!msg) return;
      const otherId = msg.sender_id === user.id ? msg.receiver_id : msg.sender_id;

      if (open && activeId === otherId) {
        setMessages((prev) => (prev.some((m) => m.id === msg.id) ? prev : [...prev, msg]));
        markChatRead(otherId).catch(() => {});
      } else {
        setUnread((prev) => ({ ...prev, [otherId]: (prev[otherId] || 0) + 1 }));
      }

      if (msg.sender_id !== user.id) playSfxRef.current("/audio/receive.mp3");
    });
    return unsub;
  }, [user, subscribe, open, activeId]);

  // Closing the modal ends the open thread, so reopening always refetches the
  // history instead of showing a thread missing the messages that arrived meanwhile.
  const setOpen = (next) => {
    const value = typeof next === "function" ? next(open) : next;
    if (!value) setActiveId(null);
    setOpenState(value);
  };

  const totalUnread = Object.values(unread).filter((count) => count > 0).length;
  const activeContact = conversations.find((c) => c.user_id === activeId) || null;

  const openChat = async (contactId) => {
    setActiveId(contactId);
    setUnread((prev) => ({ ...prev, [contactId]: 0 }));
    setMessages([]);
    setPage(1);
    setHasMore(false);
    setLoading(true);
    try {
      const res = await getDirectHistory(contactId, { page: 1 });
      const history = Array.isArray(res?.data) ? res.data : [];
      setMessages(history.slice().reverse()); // server sends newest first
      setHasMore(history.length === pageSize);
    } catch {
      /* ignore; the thread just stays empty */
    } finally {
      setLoading(false);
    }
  };

  const loadOlder = async () => {
    if (loadingOlder || !hasMore) return;
    setLoadingOlder(true);
    try {
      const nextPage = page + 1;
      const res = await getDirectHistory(activeId, { page: nextPage });
      const older = Array.isArray(res?.data) ? res.data : [];
      setMessages((prev) => [...older.slice().reverse(), ...prev]);
      setPage(nextPage);
      setHasMore(older.length === pageSize);
    } catch {
      /* ignore; keep existing messages */
    } finally {
      setLoadingOlder(false);
    }
  };

  const sendMessage = () => {
    const content = draft.trim();
    if (!content || !activeId) return;
    send("send_direct_message", { receiver_id: activeId, content });
    playSfxRef.current("/audio/send.mp3");
    setDraft("");
  };

  return {
    user,
    open,
    setOpen,
    conversations,
    activeId,
    setActiveId,
    activeContact,
    messages,
    unread,
    totalUnread,
    hasMore,
    loading,
    loadingOlder,
    draft,
    setDraft,
    openChat,
    loadOlder,
    sendMessage,
  };
}

export default function Chat({ chat }) {
  if (!chat.open) return null;

  return (
    <div
      className="chatModalOverlay"
      onClick={(e) => e.target === e.currentTarget && chat.setOpen(false)}
    >
      <div className="chatModal">
        {chat.activeContact ? (
          <SingleChat chat={chat} contact={chat.activeContact} />
        ) : (
          <Contacts chat={chat} />
        )}
      </div>
    </div>
  );
}

function Contacts({ chat }) {
  // NOTE: unfollowing a user doesn't hide contact UI, this isn't a bug
  //       contacts are not real time. Messages will not work

  return (
    <div className="chatPanel">
      <div className="chatPanelHeader">
        <h3 className="chatPanelTitle">Chat</h3>
        <button type="button" className="chatCloseButton" onClick={() => chat.setOpen(false)}>
          &times;
        </button>
      </div>
      <div className="chatPanelBody">
        {chat.conversations.length === 0 ? (
          <p className="chatEmpty">No conversations yet.</p>
        ) : (
          chat.conversations.map((c) => (
            <div
              key={c.user_id}
              role="button"
              tabIndex={0}
              className="chatContact"
              onClick={() => chat.openChat(c.user_id)}
              onKeyDown={(e) => {
                if (e.key === "Enter") chat.openChat(c.user_id);
              }}
            >
              <Avatar avatar={c.avatar} username={c.username} />
              <div className="chatContactInfo">
                <strong className="chatContactName">{displayName(c)}</strong>
              </div>
              {chat.unread[c.user_id] > 0 && (
                <span className="chatContactUnread">{chat.unread[c.user_id]}</span>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}

function SingleChat({ chat, contact }) {
  const listEndRef = useRef(null);

  // Keep the newest message in view.
  useEffect(() => {
    listEndRef.current?.scrollIntoView({ block: "end" });
  }, [chat.messages.length]);

  return (
    <div className="chatPanel chatThread">
      <div className="chatPanelHeader">
        <button
          type="button"
          className="chatBackButton"
          onClick={() => chat.setActiveId(null)}
          title="Back"
        >
          &larr;
        </button>
        <div className="chatThreadTitle">
          <Avatar avatar={contact.avatar} username={contact.username} size={32} />
          <strong className="chatContactName">{displayName(contact)}</strong>
        </div>
        <button type="button" className="chatCloseButton" onClick={() => chat.setOpen(false)}>
          &times;
        </button>
      </div>

      <div className="chatMessages">
        {chat.loading && <p className="chatEmpty">Loading...</p>}
        {!chat.loading && chat.hasMore && (
          <button
            type="button"
            className="chatLoadOlder"
            onClick={chat.loadOlder}
            disabled={chat.loadingOlder}
          >
            {chat.loadingOlder ? "Loading..." : "Load previous"}
          </button>
        )}
        {!chat.loading && chat.messages.length === 0 && (
          <p className="chatEmpty">No messages yet. Say hello!</p>
        )}
        {chat.messages.map((m) => (
          <div key={m.id} className={`chatBubble ${m.sender_id === chat.user?.id ? "mine" : "theirs"}`}>
            <span className="chatBubbleText">{m.content}</span>
            <span className="chatBubbleTime">{formatMessageTime(m.created_at)}</span>
          </div>
        ))}
        <div ref={listEndRef} />
      </div>

      <Composer draft={chat.draft} setDraft={chat.setDraft} onSend={chat.sendMessage} />
    </div>
  );
}
