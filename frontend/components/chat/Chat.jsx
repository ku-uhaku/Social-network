"use client";

import { useEffect, useRef, useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useWebSocket } from "@/contexts/WebSocketContext";
import { useAudio } from "@/contexts/AudioContext";
import { getConversations, getDirectHistory, markChatRead } from "@/lib/api/chat";
import Avatar from "@/components/shared/Avatar";
import EmojiPicker from "./EmojiPicker";
import "@/css/chat.css";

const empty_single_chat = {
  messages: [],
  loaded: false,
  unread: 0,
  hasMore: false,
  oldestId: null,
  loadingOlder: false,
};

export function useChat() {
  const { user } = useAuth();
  const { send, subscribe } = useWebSocket();
  const { playSfx } = useAudio();
  const playSfxRef = useRef(playSfx);
  useEffect(() => {
    playSfxRef.current = playSfx;
  }, [playSfx]);

  const [open, setOpen] = useState(false);
  const [conversations, setConversations] = useState([]);
  const [activeId, setActiveId] = useState(null);
  const [contacts, setContacts] = useState({}); // { [id]: { messages, loaded, unread, hasMore, oldestId } }
  const [draft, setDraft] = useState("");
  const [loading, setLoading] = useState(false);

  const getThread = (id) => contacts[id] || empty_single_chat;

  const patchThread = (id, patch) => {
    setContacts((prev) => {
      const thread = prev[id] || empty_single_chat;
      const next = typeof patch === "function" ? patch(thread) : patch;
      return { ...prev, [id]: { ...thread, ...next } };
    });
  };

  // load conversation list + seed persistent unread counts
  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    getConversations()
      .then((res) => {
        if (cancelled) return;
        const list = Array.isArray(res?.data) ? res.data : [];
        setConversations(list);
        setContacts((prev) => {
          let next = prev;
          for (const c of list) {
            if (!next[c.user_id]) {
              next = { ...next, [c.user_id]: { ...empty_single_chat, unread: c.unread_count || 0 } };
            }
          }
          return next;
        });
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [user]);

  // realtime dm
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

      if (isViewing) markChatRead(otherId).catch(() => {});

      if (msg.sender_id !== user.id) playSfxRef.current("/audio/receive.mp3");
    });
    return unsub;
  }, [user, subscribe, open, activeId]);

  const unread = Object.fromEntries(Object.entries(contacts).map(([id, c]) => [id, c.unread]));
  const totalUnread = Object.values(contacts).filter((c) => c.unread > 0).length;
  const activeContact = activeId ? conversations.find((c) => c.user_id === activeId) : null;

  const activeThread = getThread(activeId);

  const openChat = async (contactId) => {
    setActiveId(contactId);
    markChatRead(contactId).catch(() => {});
    patchThread(contactId, { unread: 0 });
    if ((contacts[contactId] || empty_single_chat).loaded) return;

    setLoading(true);
    try {
      const res = await getDirectHistory(contactId);
      const data = res?.data || {};
      const history = Array.isArray(data.messages) ? data.messages : [];
      patchThread(contactId, (prev) => {
        const seen = new Set(prev.messages.map((m) => m.id));
        const missing = history.filter((m) => !seen.has(m.id));
        return {
          messages: [...missing, ...prev.messages],
          loaded: true,
          hasMore: Boolean(data.has_more),
          oldestId: history.length ? history[0].id : prev.oldestId,
        };
      });
    } catch {
      patchThread(contactId, { loaded: false });
    } finally {
      setLoading(false);
    }
  };

  const loadOlder = async (contactId) => {
    const thread = contacts[contactId] || empty_single_chat;
    if (!thread.loaded || !thread.hasMore || thread.loadingOlder) return;
    patchThread(contactId, { loadingOlder: true });
    try {
      const res = await getDirectHistory(contactId, { beforeId: thread.oldestId });
      const data = res?.data || {};
      const older = Array.isArray(data.messages) ? data.messages : [];
      patchThread(contactId, (t) => ({
        messages: [...older, ...t.messages],
        hasMore: Boolean(data.has_more),
        oldestId: older.length ? older[0].id : t.oldestId,
        loadingOlder: false,
      }));
    } catch {
      patchThread(contactId, { loadingOlder: false });
    }
  };

  const sendMessage = () => {
    const content = draft.trim();
    if (!content || !activeId) return;
    send("send_direct_message", { receiver_id: activeId, content });
    playSfxRef.current("/audio/send.mp3");
    setDraft("");
  };

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
    activeThread,
    unread,
    totalUnread,
    draft,
    setDraft,
    loading,
    openChat,
    loadOlder,
    sendMessage,
    handleKeyDown,
    handleBackdrop,
  };
}

export default function Chat({ chat }) {
  return (
    <>
      {chat.open && (
        <div className="chatModalOverlay" onClick={chat.handleBackdrop}>
          <div className="chatModal">
            {chat.activeId && chat.activeContact ? (
              <SingleChat
                contact={chat.activeContact}
                thread={chat.activeThread}
                loading={chat.loading}
                meId={chat.user?.id}
                actions={{
                  draft: chat.draft,
                  onDraftChange: chat.setDraft,
                  onSend: chat.sendMessage,
                  onKeyDown: chat.handleKeyDown,
                  onLoadOlder: () => chat.loadOlder(chat.activeId),
                  onBack: () => chat.setActiveId(null),
                  onClose: () => chat.setOpen(false),
                }}
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
  // NOTE: unfollowing a user doesn't hide contact UI, this isn't a bug
  //       contacts are not real time. Messages will not work

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
                <strong className="chatContactName">{c.username}</strong>
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

function SingleChat({ contact, thread, loading, meId, actions }) {
  const listEndRef = useRef(null);
  const [emojiOpen, setEmojiOpen] = useState(false);

  // keep newest
  useEffect(() => {
    listEndRef.current?.scrollIntoView({ block: "end" });
  }, [thread.messages]);

  return (
    <div className="chatPanel chatThread">
      <div className="chatPanelHeader">
        <button type="button" className="chatBackButton" onClick={actions.onBack} title="Back">
          &larr;
        </button>
        <div className="chatThreadTitle">
          <Avatar avatar={contact.avatar} username={contact.username} size={32} />
          <strong className="chatContactName">{contact.username}</strong>
        </div>
        <button type="button" className="chatCloseButton" onClick={actions.onClose}>
          &times;
        </button>
      </div>

      <div className="chatMessages">
        {loading && <p className="chatEmpty">Loading...</p>}
        {!loading && thread.hasMore && (
          <button
            type="button"
            className="chatLoadOlder"
            onClick={actions.onLoadOlder}
            disabled={thread.loadingOlder}
          >
            {thread.loadingOlder ? "Loading..." : "Load previous"}
          </button>
        )}
        {!loading && thread.messages.length === 0 && (
          <p className="chatEmpty">No messages yet. Say hello!</p>
        )}
        {thread.messages.map((m) => {
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

      {emojiOpen && <EmojiPicker onPick={(e) => actions.onDraftChange(actions.draft + e)} />}

      <div className="chatComposer">
        <button
          type="button"
          className={`chatEmojiToggle ${emojiOpen ? "active" : ""}`}
          onClick={() => setEmojiOpen(!emojiOpen)}
          title="Emoji"
        >
          {"\u263A"}
        </button>
        <textarea
          className="chatInput"
          rows={1}
          value={actions.draft}
          placeholder={`Type...`}
          onChange={(e) => actions.onDraftChange(e.target.value)}
          onKeyDown={actions.onKeyDown}
        />
        <button type="button" className="chatSendButton" onClick={actions.onSend}>
          Send
        </button>
      </div>
    </div>
  );
}

function formatTime(dateString) {
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) {
    return "";
  }

  const time = date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  if (date.toDateString() === new Date().toDateString()) {
    // same time
    return time;
  }

  // if not same day add date
  return `${date.toLocaleDateString([], { month: "short", day: "numeric" })} ${time}`;
}
