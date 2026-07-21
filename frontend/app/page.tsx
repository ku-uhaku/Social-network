// app/page.tsx
"use client";

import React, { useState } from "react";
import HollowKnightButton from "@/components/HollowKnightButton";
import HollowKnightInput from "@/components/HollowKnightInput";
import HollowKnightSelect from "@/components/HollowKnightSelect";
import HollowKnightCheckbox from "@/components/HollowKnightCheckbox";
import HollowKnightRadio from "@/components/HollowKnightRadio";

export default function Home() {
  const [selectedCharm, setSelectedCharm] = useState("compass");

  return (
    <main className="min-h-screen bg-hk-abyss flex items-center justify-center p-6 text-hk-bone selection:bg-hk-soul/30 selection:text-hk-bone">
      <div className="w-full max-w-md bg-hk-surface/40 border border-hk-border rounded-xl p-8 backdrop-blur-md shadow-2xl flex flex-col gap-6">
        
        {/* Header Title */}
        <div className="text-center">
          <span className="text-hk-muted text-[10px] tracking-[0.3em] uppercase block mb-1">
            ~ Hallownest ~
          </span>
          <h1 className="text-2xl font-serif tracking-[0.2em] text-hk-bone uppercase [text-shadow:0_0_12px_rgba(226,232,240,0.3)]">
            Knight Profile
          </h1>
          <div className="w-28 h-[1px] bg-gradient-to-r from-transparent via-hk-border to-transparent my-2 mx-auto" />
        </div>

        {/* Text Input */}
        <HollowKnightInput label="Vessel Name" placeholder="e.g. Ghost" />

        {/* Dropdown Select */}
        <HollowKnightSelect
          label="Nail Upgrade Level"
          options={[
            { value: "old", label: "Old Nail" },
            { value: "sharpened", label: "Sharpened Nail" },
            { value: "chiseled", label: "Chiseled Nail" },
            { value: "pure", label: "Pure Nail" },
          ]}
        />

        {/* Checkbox Group */}
        <div className="flex flex-col gap-2 pt-1">
          <label className="text-xs uppercase tracking-[0.2em] text-hk-muted font-serif">
            Abilities
          </label>
          <HollowKnightCheckbox label="Mantis Claw" defaultChecked />
          <HollowKnightCheckbox label="Monarch Wings" />
        </div>

        {/* Radio Group */}
        <div className="flex flex-col gap-2 pt-1">
          <label className="text-xs uppercase tracking-[0.2em] text-hk-muted font-serif">
            Equipped Charm
          </label>
          <HollowKnightRadio
            name="charm"
            label="Wayward Compass"
            checked={selectedCharm === "compass"}
            onChange={() => setSelectedCharm("compass")}
          />
          <HollowKnightRadio
            name="charm"
            label="Gathering Swarm"
            checked={selectedCharm === "swarm"}
            onChange={() => setSelectedCharm("swarm")}
          />
        </div>

        {/* Submit Action Button */}
        <div className="pt-4 flex justify-center">
          <HollowKnightButton>Save Settings</HollowKnightButton>
        </div>

      </div>
    </main>
  );
}