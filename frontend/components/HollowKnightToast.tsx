// components/HollowKnightToast.tsx
"use client";

import React, { useEffect } from "react";
import { useAudio } from "@/contexts/AudioContext";

interface ToastProps {
  message: string;
  subtitle?: string;
  isVisible: boolean;
  onClose: () => void;
  duration?: number;
}

export default function HollowKnightToast({
  message,
  subtitle,
  isVisible,
  onClose,
  duration = 3500,
}: ToastProps) {
  const { playEffect } = useAudio();

  useEffect(() => {
    if (isVisible) {
      playEffect("/audio/toast.mp3");

      const timer = setTimeout(() => {
        onClose();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [isVisible, duration, onClose]);

  if (!isVisible) return null;

  return (
    <div className="fixed bottom-8 right-8 z-50 transition-all duration-300 transform animate-in slide-in-from-bottom-4 fade-in">
      <div className="relative flex items-center gap-4 bg-[#090d12]/95 border border-[#2d3f52] rounded-sm px-6 py-3.5 shadow-[0_8px_30px_rgba(0,0,0,0.85)] backdrop-blur-md min-w-[280px]">
        
        {/* Inner Border Frame */}
        <div className="absolute inset-1 border border-[#16202c] pointer-events-none rounded-sm" />

        {/* Soul Glow Indicator Icon */}
        <div className="flex items-center justify-center text-[#38bdf8] text-sm select-none [text-shadow:0_0_10px_rgba(56,189,248,0.7)]">
          ❖
        </div>

        {/* Notification Text */}
        <div className="flex flex-col flex-1 pr-2">
          <span className="font-serif text-xs tracking-[0.2em] text-[#e2e8f0] uppercase [text-shadow:0_0_8px_rgba(255,255,255,0.15)]">
            {message}
          </span>
          {subtitle && (
            <span className="font-serif text-[11px] tracking-wider text-[#64748b] mt-0.5">
              {subtitle}
            </span>
          )}
        </div>

        {/* Close Button */}
        <button
          onClick={onClose}
          className="relative z-10 text-[#64748b] hover:text-[#e2e8f0] text-xs transition-colors duration-200 cursor-pointer p-1"
        >
          ✕
        </button>
      </div>
    </div>
  );
}