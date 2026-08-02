"use client";

import { useState } from "react";
import HollowKnightInput from "@/components/HollowKnightInput";
import HollowKnightCheckbox from "@/components/HollowKnightCheckbox";
import HollowKnightSelect from "@/components/HollowKnightSelect";
import HollowKnightButton from "@/components/HollowKnightButton";
import HollowKnightModal from "@/components/HollowKnightModal";
import HollowKnightToast from "@/components/HollowKnightToast";
import { useAudio, TRACKS } from "@/contexts/AudioContext";

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

  return (
    <main className="page">
      <div className="panel">
        <div className="panel-header">
          <h2>Knight Profile</h2>
          <div className="panel-divider" />
        </div>

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
          onChange={(e) => setCurrentTrack(e.target.value)}
          options={TRACKS.map((t) => ({ value: t.id, label: t.title }))}
        />

        <div className="field-group">
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

        <div className="actions">
          <HollowKnightButton onClick={() => setIsToastVisible(true)}>
            Save Settings
          </HollowKnightButton>
          <HollowKnightButton onClick={() => setIsModalOpen(true)}>
            View Charm Info
          </HollowKnightButton>
        </div>
      </div>

      <HollowKnightModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title="Charm Acquired"
      >
        <p>
          You have acquired the <strong>Wayward Compass</strong>. It points
          toward your current location when resting at benches or traveling
          through Hallownest.
        </p>
      </HollowKnightModal>

      <HollowKnightToast
        isVisible={isToastVisible}
        message="Settings Saved"
        subtitle="Your audio preferences have been updated."
        onClose={() => setIsToastVisible(false)}
      />
    </main>
  );
}
