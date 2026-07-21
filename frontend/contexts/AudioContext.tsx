// contexts/AudioContext.tsx
"use client";

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  ReactNode,
} from "react";

export type TrackId = "dirtmouth" | "greenpath" | "reflection";

interface Track {
  id: TrackId;
  title: string;
  src: string;
}

export const TRACKS: Track[] = [
  { id: "dirtmouth", title: "Dirtmouth", src: "/audio/02. Dirtmouth.mp3" },
  { id: "greenpath", title: "Greenpath", src: "/audio/05. Greenpath.mp3" },
  { id: "reflection", title: "Reflection", src: "/audio/07. Reflection.mp3" },
];

const DEFAULT_TRACK: TrackId = "greenpath";

interface AudioContextValue {
  currentTrack: TrackId;
  setCurrentTrack: (id: TrackId) => void;
  isMusicMuted: boolean;
  toggleMusicMuted: () => void;
  setMusicMuted: (muted: boolean) => void;
  isEffectsMuted: boolean;
  toggleEffectsMuted: () => void;
  setEffectsMuted: (muted: boolean) => void;
  musicVolume: number;
  setMusicVolume: (v: number) => void;
  playEffect: (src: string) => void;
}

const AudioCtx = createContext<AudioContextValue | undefined>(undefined);

export function AudioProvider({ children }: { children: ReactNode }) {
  const musicRef = useRef<HTMLAudioElement | null>(null);
  const ambientRef = useRef<HTMLAudioElement | null>(null);

  const [currentTrack, setCurrentTrack] = useState<TrackId>(DEFAULT_TRACK);
  const [isMusicMuted, setIsMusicMuted] = useState(false);
  const [isEffectsMuted, setIsEffectsMuted] = useState(false);
  const [musicVolume, setMusicVolume] = useState(0.5);

  // 1. Initialize Music & Ambient Wind Elements
  useEffect(() => {
    // Music element
    const music = new Audio();
    music.loop = true;
    musicRef.current = music;

    // Ambient loop (Cave Wind)
    const ambient = new Audio("/audio/Cave Wind Loop.mp3");
    ambient.loop = true;
    ambient.volume = 0.3; // Gentle wind background
    ambientRef.current = ambient;

    return () => {
      music.pause();
      music.src = "";
      ambient.pause();
      ambient.src = "";
    };
  }, []);

  // 2. Change Music Track
  useEffect(() => {
    const music = musicRef.current;
    const track = TRACKS.find((t) => t.id === currentTrack);
    if (!music || !track) return;

    music.src = track.src;
    music.volume = musicVolume;
    music.muted = isMusicMuted;

    music.play().catch(() => {});
  }, [currentTrack]);

  // 3. Play/Pause Ambient Loop based on isEffectsMuted
  useEffect(() => {
    const ambient = ambientRef.current;
    if (!ambient) return;

    ambient.muted = isEffectsMuted;

    if (!isEffectsMuted) {
      ambient.play().catch(() => {});
    } else {
      ambient.pause();
    }
  }, [isEffectsMuted]);

  // 4. Resume both on first user interaction if blocked
  useEffect(() => {
    const resume = () => {
      if (musicRef.current && !isMusicMuted) {
        musicRef.current.play().catch(() => {});
      }
      if (ambientRef.current && !isEffectsMuted) {
        ambientRef.current.play().catch(() => {});
      }
      window.removeEventListener("click", resume);
      window.removeEventListener("keydown", resume);
    };

    window.addEventListener("click", resume);
    window.addEventListener("keydown", resume);

    return () => {
      window.removeEventListener("click", resume);
      window.removeEventListener("keydown", resume);
    };
  }, [isMusicMuted, isEffectsMuted]);

  // Sync Music Mute
  useEffect(() => {
    if (musicRef.current) musicRef.current.muted = isMusicMuted;
  }, [isMusicMuted]);

  // Sync Volume
  useEffect(() => {
    if (musicRef.current) musicRef.current.volume = musicVolume;
  }, [musicVolume]);

  const toggleMusicMuted = () => setIsMusicMuted((m) => !m);
  const toggleEffectsMuted = () => setIsEffectsMuted((m) => !m);

  // Play temporary SFX on user actions (e.g., hover/click)
  const playEffect = (src: string) => {
    if (isEffectsMuted) return;
    const sfx = new Audio(src);
    sfx.volume = musicVolume;
    sfx.play().catch(() => {});
  };

  return (
    <AudioCtx.Provider
      value={{
        currentTrack,
        setCurrentTrack,
        isMusicMuted,
        toggleMusicMuted,
        setMusicMuted: setIsMusicMuted,
        isEffectsMuted,
        toggleEffectsMuted,
        setEffectsMuted: setIsEffectsMuted,
        musicVolume,
        setMusicVolume,
        playEffect,
      }}
    >
      {children}
    </AudioCtx.Provider>
  );
}

export function useAudio() {
  const ctx = useContext(AudioCtx);
  if (!ctx) throw new Error("useAudio must be used within an AudioProvider");
  return ctx;
}