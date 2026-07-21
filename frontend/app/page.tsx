// app/page.tsx
"use client";

import React, { useState } from "react";
import HollowKnightInput from "@/components/HollowKnightInput";
import HollowKnightCheckbox from "@/components/HollowKnightCheckbox";
import HollowKnightSelect from "@/components/HollowKnightSelect";
import HollowKnightButton from "@/components/HollowKnightButton";
import HollowKnightModal from "@/components/HollowKnightModal";
import HollowKnightToast from "@/components/HollowKnightToast";
import { useAudio, TRACKS, TrackId } from "@/contexts/AudioContext";

export default function SettingsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isToastVisible, setIsToastVisible] = useState(false);

  const {
    currentTrack,
    setCurrentTrack,
    isMusicMuted,
    toggleMusicMuted,
    isEffectsMuted,
    toggleEffectsMuted,
  } = useAudio();

  const handleSave = () => {
    setIsToastVisible(true);
  };

  return (
    <main className="min-h-screen bg-hk-abyss flex flex-col items-center justify-center p-6 text-hk-bone">
      <div className="w-full max-w-md bg-hk-surface/40 border border-hk-border rounded-xl p-8 backdrop-blur-sm shadow-2xl flex flex-col gap-6">
        
        {/* Header */}
        <div className="text-center">
          <h2 className="text-2xl font-serif tracking-[0.2em] text-hk-bone uppercase">
            Knight Profile
          </h2>
          <div className="w-24 h-[1px] bg-gradient-to-r from-transparent via-hk-border to-transparent my-2 mx-auto" />
        </div>

        {/* Form Controls */}
        <HollowKnightInput label="Vessel Name" placeholder="e.g. Ghost" />

        <HollowKnightSelect
          label="Nail Upgrade Level"
          options={[
            { value: "old", label: "Old Nail" },
            { value: "sharpened", label: "Sharpened Nail" },
            { value: "chiseled", label: "Chiseled Nail" },
            { value: "pure", label: "Pure Nail" },
          ]}
        />

        <HollowKnightSelect
          label="Background Track"
          value={currentTrack}
          onChange={(e) => setCurrentTrack(e.target.value as TrackId)}
          options={TRACKS.map((t) => ({ value: t.id, label: t.title }))}
        />

        <div className="flex flex-col gap-2 pt-2">
          <HollowKnightCheckbox label="Enable Wayward Compass" />
          <HollowKnightCheckbox label="Show Infection Glow" />
          <HollowKnightCheckbox
            label="Mute Music"
            checked={isMusicMuted}
            onChange={toggleMusicMuted}
          />
          <HollowKnightCheckbox
            label="Mute Sound Effects"
            checked={isEffectsMuted}
            onChange={toggleEffectsMuted}
          />
        </div>

        {/* Action Buttons */}
        <div className="pt-4 flex flex-col items-center gap-3">
          <HollowKnightButton onClick={handleSave}>
            Save Settings
          </HollowKnightButton>

          <HollowKnightButton onClick={() => setIsModalOpen(true)}>
            View Charm Info
          </HollowKnightButton>
        </div>

      </div>

      {/* Hallownest Modal */}
      <HollowKnightModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Charm Acquired"
      >
        <p>
          You have acquired the <strong>Wayward Compass</strong>. It points toward
          your current location when resting at benches or traveling through Hallownest.
        </p>
      </HollowKnightModal>

      {/* Hallownest Toast */}
      <HollowKnightToast
        isVisible={isToastVisible}
        message="Settings Saved"
        subtitle="Your audio preferences have been updated."
        onClose={() => setIsToastVisible(false)}
      />
    </main>
  );
}