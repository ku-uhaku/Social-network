// components/HollowKnightRadio.tsx
"use client";

import React from "react";
import Image from "next/image";

interface HKRadioProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
}

export default function HollowKnightRadio({ label, className = "", ...props }: HKRadioProps) {
  return (
    <label className="group relative inline-flex items-center gap-3 cursor-pointer select-none py-1">
      {/* Hidden real radio */}
      <input type="radio" className="sr-only peer" {...props} />

      {/* Left Pointer on Hover/Focus */}
      <span className="opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 text-hk-bone text-xs">
        ‹<span className="text-[9px] -ml-0.5">❖</span>
      </span>

      {/* Custom Radio Circle — rounded-full frame */}
      <div
        className="
          relative w-5 h-5 bg-hk-surface border border-hk-border rounded-full
          flex items-center justify-center transition-all duration-300
          peer-checked:border-hk-soul peer-checked:bg-hk-soul/10
          peer-checked:[box-shadow:0_0_10px_rgba(0,229,255,0.3)]
          group-hover:border-hk-muted overflow-hidden
          peer-checked:[&>span]:opacity-100 peer-checked:[&>span]:scale-100
        "
      >
        {/* Nested span — styled via the parent's arbitrary variant above */}
        <span className="relative w-full h-full p-0.5 opacity-0 scale-0 transition-all duration-200 [filter:drop-shadow(0_0_4px_rgba(0,229,255,0.8))]">
          <Image
            src="/hollow.png"
            alt=""
            fill
            sizes="20px"
            className="object-contain"
          />
        </span>
      </div>

      {/* Label Text */}
      <span className="font-serif tracking-widest text-sm text-hk-muted group-hover:text-hk-bone group-hover:[text-shadow:0_0_8px_rgba(255,255,255,0.4)] transition-all duration-300">
        {label}
      </span>
    </label>
  );
}