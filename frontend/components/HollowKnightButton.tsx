"use client";

import React from "react";

interface HKButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
}

export default function HollowKnightButton({ children, className = "", ...props }: HKButtonProps) {
  return (
    <button
      className={`group relative inline-flex items-center justify-center gap-3 py-2 px-6 bg-transparent cursor-pointer transition-all duration-300 select-none ${className}`}
      {...props}
    >
      {/* Left Menu Pointer Ornament */}
      <span className="opacity-0 -translate-x-3 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs flex items-center gap-0.5">
        ‹<span className="text-[10px] -ml-1">❖</span>
      </span>

      {/* Button Text */}
      <span className="font-serif tracking-[0.2em] uppercase text-hk-muted group-hover:text-hk-bone group-hover:scale-105 transition-all duration-300 group-hover:[text-shadow:0_0_12px_rgba(255,255,255,0.7),0_0_20px_rgba(0,229,255,0.4)]">
        {children}
      </span>

      {/* Right Menu Pointer Ornament */}
      <span className="opacity-0 translate-x-3 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs flex items-center gap-0.5">
        <span className="text-[10px] -mr-1">❖</span>›
      </span>
    </button>
  );
}