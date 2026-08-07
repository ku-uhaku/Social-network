"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";
import { getConversations, getDirectHistory } from "@/lib/api/chat";
import Avatar from "@/components/shared/Avatar";
import "@/css/chat.css";

const empty_single_chat = { messages: [], loaded: false, unread: 0 };

function useChat() {
  const { user } = useAuth();
  const { send, subscribe } = useWebSocket();

  const [open, setOpen] = useState(false);
  const [conversations, setConversations] = useState([]);
  const [activeId, setActiveId] = useState(null);
  const [contacts, setContacts] = useState({}); // { [id]: { messages, loaded, unread } }
  const [draft, setDraft] = useState("");
  const [loading, setLoading] = useState(false);

  const getThread = (id) => contacts[id] || empty_single_chat;

  const patchThread = useCallback((id, patch) => {
    setContacts((prev) => {
      const thread = prev[id] || empty_single_chat;
      const next = typeof patch === "function" ? patch(thread) : patch;
      return { ...prev, [id]: { ...thread, ...next } };
    });
  }, [contacts]);

  // load conversation list
  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    getConversations()
      .then((res) => !cancelled && setConversations(Array.isArray(res?.data) ? res.data : []))
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [user]);

  // dm
  useEffect(() => {
    if (!user) return;
    const unsub = subscribe("new_direct_message", (msg) => {
      if (!msg) return;
      const otherId = msg.sender_id === user.id ? msg.receiver_id : msg.sender_id;
      const isViewing = open && activeId === otherId;

      setContacts((prev) => {
        const thread = prev[otherId] || empty_single_chat;
        if (thread.messages.some((m) => m.id === msg.id)) return prev;
        return {
          ...prev,
          [otherId]: {
            ...thread,
            messages: [...thread.messages, msg],
            unread: isViewing ? thread.unread : thread.unread + 1,
          },
        };
      });
    });
    return unsub;
  }, [user, subscribe, open, activeId]);

  const unread = useMemo(
    () => Object.fromEntries(Object.entries(contacts).map(([id, c]) => [id, c.unread])),
    [contacts]
  );

  const totalUnread = useMemo(
    () => Object.values(contacts).reduce((sum, c) => sum + c.unread, 0),
    [contacts]
  );

  const activeContact = useMemo(
    () => (activeId ? conversations.find((c) => c.user_id === activeId) : null),
    [activeId, conversations]
  );
  const threadMessages = getThread(activeId).messages;

  const openChat = useCallback(
    async (contactId) => {
      setActiveId(contactId);
      patchThread(contactId, { unread: 0 });
      if (getThread(contactId).loaded) return;

      setLoading(true);
      try {
        const res = await getDirectHistory(contactId);
        const history = Array.isArray(res?.data) ? res.data : [];
        patchThread(contactId, (thread) => {
          const missing = history.filter((m) => !thread.messages.some((cm) => cm.id === m.id));
          return { messages: [...missing, ...thread.messages], loaded: true };
        });
      } catch {
        patchThread(contactId, { loaded: false });
      } finally {
        setLoading(false);
      }
    },
    [patchThread]
  );

  const sendMessage = useCallback(() => {
    const content = draft.trim();
    if (!content || !activeId) return;
    send("send_direct_message", { receiver_id: activeId, content });
    setDraft("");
  }, [draft, activeId, send]);

  const handleKeyDown = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const handleBackdrop = (e) => {
    if (e.target === e.currentTarget) setOpen(false);
  };

  return {
    user,
    open,
    setOpen,
    conversations,
    activeId,
    setActiveId,
    activeContact,
    threadMessages,
    unread,
    totalUnread,
    draft,
    setDraft,
    loading,
    openChat,
    sendMessage,
    handleKeyDown,
    handleBackdrop,
  };
}

export default function ChatDock() {
  const chat = useChat();

  return (
    <>
      <aside className="chatDock">
        <button
          type="button"
          className={`chatDockButton ${chat.open ? "active" : ""}`}
          title="Chat"
          onClick={() => chat.setOpen(!chat.open)}
        >
          <img className="chatDockIcon" src="/images/icon_chat.png" alt="Chat" />
          {chat.totalUnread > 0 && <span className="chatDockBadge">{chat.totalUnread}</span>}
        </button>
        {/* Groups button will be added here in phase 2 */}
      </aside>

      {chat.open && (
        <div className="chatModalOverlay" onClick={chat.handleBackdrop}>
          <div className="chatModal">
            {chat.activeId && chat.activeContact ? (
              <SingleChat
                contact={chat.activeContact}
                messages={chat.threadMessages}
                draft={chat.draft}
                loading={chat.loading}
                meId={chat.user?.id}
                onDraftChange={chat.setDraft}
                onSend={chat.sendMessage}
                onKeyDown={chat.handleKeyDown}
                onBack={() => chat.setActiveId(null)}
                onClose={() => chat.setOpen(false)}
              />
            ) : (
              <Contacts
                conversations={chat.conversations}
                unread={chat.unread}
                loading={chat.loading}
                onOpen={chat.openChat}
                onClose={() => chat.setOpen(false)}
              />
            )}
          </div>
        </div>
      )}
    </>
  );
}

function Contacts({ conversations, unread, loading, onOpen, onClose }) {
  return (
    <div className="chatPanel">
      <div className="chatPanelHeader">
        <h3 className="chatPanelTitle">Chat</h3>
        <button type="button" className="chatCloseButton" onClick={onClose}>
          &times;
        </button>
      </div>
      <div className="chatPanelBody">
        {loading ? (
          <p className="chatEmpty">Loading...</p>
        ) : conversations.length === 0 ? (
          <p className="chatEmpty">No conversations yet.</p>
        ) : (
          conversations.map((c) => (
            <div
              key={c.user_id}
              role="button"
              tabIndex={0}
              className="chatContact"
              onClick={() => onOpen(c.user_id)}
              onKeyDown={(e) => {
                if (e.key === "Enter") onOpen(c.user_id);
              }}
            >
              <Avatar avatar={c.avatar} username={c.username} />
              <div className="chatContactInfo">
                <strong className="chatContactName">@{c.username}</strong>
              </div>
              {unread[c.user_id] > 0 && (
                <span className="chatContactUnread">{unread[c.user_id]}</span>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}

function SingleChat({
  contact,
  messages,
  draft,
  loading,
  meId,
  onDraftChange,
  onSend,
  onKeyDown,
  onBack,
  onClose,
}) {
  const listEndRef = useRef(null);

  // keep newest
  useEffect(() => {
    listEndRef.current?.scrollIntoView({ block: "end" });
  }, [messages]);

  return (
    <div className="chatPanel chatThread">
      <div className="chatPanelHeader">
        <button type="button" className="chatBackButton" onClick={onBack} title="Back">
          &larr;
        </button>
        <div className="chatThreadTitle">
          <Avatar avatar={contact.avatar} username={contact.username} size={32} />
          <strong className="chatContactName">@{contact.username}</strong>
        </div>
        <button type="button" className="chatCloseButton" onClick={onClose}>
          &times;
        </button>
      </div>

      <div className="chatMessages">
        {loading && <p className="chatEmpty">Loading...</p>}
        {!loading && messages.length === 0 && (
          <p className="chatEmpty">No messages yet. Say hello!</p>
        )}
        {messages.map((m) => {
          const mine = m.sender_id === meId;
          return (
            <div key={m.id} className={`chatBubble ${mine ? "mine" : "theirs"}`}>
              <span className="chatBubbleText">{m.content}</span>
              <span className="chatBubbleTime">{formatTime(m.created_at)}</span>
            </div>
          );
        })}
        <div ref={listEndRef} />
      </div>

      <div className="chatComposer">
        <textarea
          className="chatInput"
          rows={1}
          value={draft}
          placeholder={`Message @${contact.username}...`}
          onChange={(e) => onDraftChange(e.target.value)}
          onKeyDown={onKeyDown}
        />
        <button type="button" className="chatSendButton" onClick={onSend}>
          Send
        </button>
      </div>
    </div>
  );
}

function formatTime(dateString) {
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())){
    return "";
  }

  const time = date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  if (date.toDateString() === new Date().toDateString()){
    // same time
    return time;
  }

  // if not same day add date
  return `${date.toLocaleDateString([], { month: "short", day: "numeric" })} ${time}`;
}