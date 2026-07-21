"use client";

import React from "react";

interface HKSelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  options: { value: string; label: string }[];
}

export default function HollowKnightSelect({ label, options, className = "", ...props }: HKSelectProps) {
  return (
    <div className="flex flex-col gap-1.5 w-full max-w-sm">
      {label && (
        <label className="text-xs uppercase tracking-[0.2em] text-hk-muted font-serif">
          {label}
        </label>
      )}
      <div className="relative flex items-center group">
        {/* Left Pointer */}
        <span className="absolute -left-5 opacity-0 -translate-x-2 group-focus-within:opacity-100 group-focus-within:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
          ‹<span className="text-[9px] -ml-0.5">❖</span>
        </span>

        <select
          className={`w-full bg-hk-surface border border-hk-border rounded px-4 py-2 pr-8 text-hk-bone font-serif tracking-wider appearance-none focus:outline-none focus:border-hk-soul focus:[box-shadow:0_0_12px_rgba(0,229,255,0.25)] cursor-pointer transition-all duration-300 ${className}`}
          {...props}
        >
          {options.map((opt) => (
            <option key={opt.value} value={opt.value} className="bg-hk-surface text-hk-bone">
              {opt.label}
            </option>
          ))}
        </select>

        {/* Down Arrow Ornament */}
        <div className="absolute right-3 pointer-events-none text-hk-muted text-xs group-hover:text-hk-bone transition-colors">
          ❖
        </div>

        {/* Right Pointer */}
        <span className="absolute -right-5 opacity-0 translate-x-2 group-focus-within:opacity-100 group-focus-within:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
          <span className="text-[9px] -mr-0.5">❖</span>›
        </span>
      </div>
    </div>
  );
}