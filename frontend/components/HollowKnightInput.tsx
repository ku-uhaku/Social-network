// components/HollowKnightInput.tsx
"use client";

import React from "react";

interface HKInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
}

export default function HollowKnightInput({ label, className = "", ...props }: HKInputProps) {
  return (
    <div className="flex flex-col gap-1.5 w-full max-w-sm">
      {label && (
        <label className="text-xs uppercase tracking-[0.2em] text-hk-muted font-serif">
          {label}
        </label>
      )}
      <div className="relative flex items-center group">
        {/* Focus Flourish - Left Pointer */}
        <span className="absolute -left-5 opacity-0 -translate-x-2 group-focus-within:opacity-100 group-focus-within:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
          ‹<span className="text-[9px] -ml-0.5">❖</span>
        </span>

        <input
          className={`w-full bg-hk-surface/80 border border-hk-border rounded px-4 py-2 text-hk-bone font-serif tracking-wider placeholder:text-hk-muted/50 focus:outline-none focus:border-hk-soul focus:[box-shadow:0_0_12px_rgba(0,229,255,0.25)] transition-all duration-300 ${className}`}
          {...props}
        />

        {/* Focus Flourish - Right Pointer */}
        <span className="absolute -right-5 opacity-0 translate-x-2 group-focus-within:opacity-100 group-focus-within:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
          <span className="text-[9px] -mr-0.5">❖</span>›
        </span>
      </div>
    </div>
  );
}