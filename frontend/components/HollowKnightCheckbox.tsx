"use client";

import React from "react";

interface HKCheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
}

export default function HollowKnightCheckbox({ label, className = "", ...props }: HKCheckboxProps) {
  return (
    <label className="group relative inline-flex items-center gap-3 cursor-pointer select-none py-1">
      {/* Hidden real checkbox */}
      <input type="checkbox" className="sr-only peer" {...props} />

      {/* Left Pointer on Hover/Focus */}
      <span className="opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
        ‹<span className="text-[9px] -ml-0.5">❖</span>
      </span>

      {/* Custom Checkbox Box */}
      <div className="w-5 h-5 bg-hk-surface border border-hk-border rounded-sm flex items-center justify-center transition-all duration-300 peer-checked:border-hk-soul peer-checked:bg-hk-soul/10 peer-checked:[box-shadow:0_0_10px_rgba(0,229,255,0.3)] group-hover:border-hk-muted">
        {/* Custom Diamond Inner Icon */}
        <div className="w-2.5 h-2.5 bg-hk-soul rounded-[1px] rotate-45 scale-0 peer-checked:scale-100 transition-transform duration-200 [box-shadow:0_0_6px_rgba(0,229,255,0.8)]" />
      </div>

      {/* Label Text */}
      <span className="font-serif tracking-widest text-sm text-hk-muted group-hover:text-hk-bone group-hover:[text-shadow:0_0_8px_rgba(255,255,255,0.4)] transition-all duration-300">
        {label}
      </span>
    </label>
  );
}