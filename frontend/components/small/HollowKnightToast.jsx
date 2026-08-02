"use client";

import { useEffect } from "react";
import { useAudio } from "@/contexts/AudioContext";

export default function HollowKnightToast({
  message,
  subtitle,
  isVisible,
  onClose,
  duration = 3500,
}) {
  const { playEffect } = useAudio();

  useEffect(() => {
    if (!isVisible) return;
    playEffect("/audio/toast.mp3");
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [isVisible, duration, onClose]);

  if (!isVisible) return null;

  return (
    <div className="hk-toast">
      <div className="hk-toast-text">
        <span className="hk-toast-message">{message}</span>
        {subtitle && <span className="hk-toast-subtitle">{subtitle}</span>}
      </div>
      <button className="hk-toast-close" onClick={onClose}>
        ✕
      </button>
    </div>
  );
}
