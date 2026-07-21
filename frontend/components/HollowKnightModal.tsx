// components/HollowKnightModal.tsx
"use client";

import React, { useEffect } from "react";
import { useAudio } from "@/contexts/AudioContext";

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
}

export default function HollowKnightModal({
  isOpen,
  onClose,
  title,
  children,
}: ModalProps) {
  const { playEffect } = useAudio();

  useEffect(() => {
    if (isOpen) {
      playEffect("/audio/model.mp3");
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Dark overlay backdrop with fade-in */}
      <div
        className="fixed inset-0 bg-[#05080c]/85 backdrop-blur-md transition-opacity duration-300"
        onClick={onClose}
      />

      {/* Main Container */}
      <div className="relative w-full max-w-lg bg-[#0b0f14] border-2 border-[#2e3e4e] rounded-sm p-7 shadow-[0_0_35px_rgba(0,0,0,0.9)] z-10 text-hk-bone transition-all transform animate-in fade-in zoom-in-95 duration-200">
        
        {/* Inner double border outline */}
        <div className="absolute inset-1.5 border border-[#1b2633] pointer-events-none rounded-sm" />

        {/* Top Emblem Header */}
        <div className="text-center relative pt-2 pb-1">
          <div className="flex items-center justify-center gap-3">
            <span className="h-[1px] w-12 bg-gradient-to-r from-transparent to-[#4a637d]" />
            <span className="text-[#a8c0d6] text-xs select-none">❖</span>
            <span className="h-[1px] w-12 bg-gradient-to-l from-transparent to-[#4a637d]" />
          </div>

          <h2 className="text-xl font-serif tracking-[0.25em] text-[#e2e8f0] uppercase mt-2 [text-shadow:0_0_10px_rgba(226,232,240,0.2)]">
            {title}
          </h2>
        </div>

        {/* Modal Body */}
        <div className="relative my-5 px-2 font-serif text-sm leading-relaxed text-[#94a3b8]">
          {children}
        </div>

        {/* Modal Footer */}
        <div className="relative mt-6 flex justify-center pt-2">
          <button
            onClick={onClose}
            className="group relative inline-flex items-center gap-2 py-1.5 px-6 bg-[#131b24] hover:bg-[#1c2836] border border-[#2e3e4e] hover:border-[#4a637d] rounded-sm cursor-pointer font-serif text-xs tracking-[0.2em] uppercase text-[#cbd5e1] hover:text-white transition-all duration-300 shadow-sm"
          >
            <span className="opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
              ‹
            </span>
            Close
            <span className="opacity-0 translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
              ›
            </span>
          </button>
        </div>
      </div>
    </div>
  );
}