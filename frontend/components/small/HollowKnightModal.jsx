"use client";

import { useEffect } from "react";
import { useAudio } from "@/contexts/AudioContext";

export default function HollowKnightModal({ isOpen, onClose, title, children }) {
  const { playEffect } = useAudio();

  useEffect(() => {
    if (isOpen) playEffect("/audio/model.mp3");
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="hk-modal-backdrop" onClick={onClose}>
      <div className="hk-modal" onClick={(e) => e.stopPropagation()}>
        <h2>{title}</h2>
        <div className="hk-modal-body">{children}</div>
        <div className="hk-modal-footer">
          <button className="hk-button" onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
