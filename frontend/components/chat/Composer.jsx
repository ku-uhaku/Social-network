"use client";

import { useState } from "react";
import EmojiPicker from "./EmojiPicker";

// Shared message input (emoji tray + textarea + send) for direct and group chats.
export default function Composer({ draft, setDraft, onSend }) {
  const [emojiOpen, setEmojiOpen] = useState(false);

  const handleKeyDown = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      onSend();
    }
  };

  return (
    <>
      {emojiOpen && <EmojiPicker onPick={(e) => setDraft(draft + e)} />}

      <div className="chatComposer">
        <button
          type="button"
          className={`chatEmojiToggle ${emojiOpen ? "active" : ""}`}
          onClick={() => setEmojiOpen(!emojiOpen)}
          title="Emoji"
        >
          {"☺"}
        </button>
        <textarea
          className="chatInput"
          rows={1}
          value={draft}
          placeholder="Type..."
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={handleKeyDown}
          maxLength={2000}
        />
        <button type="button" className="chatSendButton" onClick={onSend}>
          Send
        </button>
      </div>
    </>
  );
}
