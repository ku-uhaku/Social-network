// app/page.tsx
"use client";

import HollowKnightInput from "@/components/HollowKnightInput";
import HollowKnightCheckbox from "@/components/HollowKnightCheckbox";
import HollowKnightSelect from "@/components/HollowKnightSelect";
import HollowKnightButton from "@/components/HollowKnightButton";
import { useAudio, TRACKS, TrackId } from "@/contexts/AudioContext";

export default function SettingsPage() {
  const {
    currentTrack,
    setCurrentTrack,
    isMusicMuted,
    toggleMusicMuted,
    isEffectsMuted,
    toggleEffectsMuted,
  } = useAudio();

  return (
    <main className="min-h-screen bg-hk-abyss flex items-center justify-center p-6 text-hk-bone">
      <div className="w-full max-w-md bg-hk-surface/40 border border-hk-border rounded-xl p-8 backdrop-blur-sm shadow-2xl flex flex-col gap-6">
        <div className="text-center">
          <h2 className="text-2xl font-serif tracking-[0.2em] text-hk-bone uppercase">
            Knight Profile
          </h2>
          <div className="w-24 h-[1px] bg-gradient-to-r from-transparent via-hk-border to-transparent my-2 mx-auto" />
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

        <div className="pt-4 flex justify-center">
          <HollowKnightButton>Save Settings</HollowKnightButton>
        </div>
      </div>
    </main>
  );
}