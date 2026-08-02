"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";

export const TRACKS = [
  { id: "dirtmouth", title: "Dirtmouth", src: "/audio/02. Dirtmouth.mp3" },
  { id: "greenpath", title: "Greenpath", src: "/audio/05. Greenpath.mp3" },
  { id: "reflection", title: "Reflection", src: "/audio/07. Reflection.mp3" },
];

const AudioCtx = createContext(null);

export function AudioProvider({ children }) {
  const musicRef = useRef(null);
  const ambientRef = useRef(null);

  const [currentTrack, setCurrentTrack] = useState("greenpath");
  const [isMusicMuted, setIsMusicMuted] = useState(false);
  const [isEffectsMuted, setIsEffectsMuted] = useState(false);

  // Create the music + ambient wind elements once.
  useEffect(() => {
    const music = new Audio();
    music.loop = true;
    musicRef.current = music;

    const ambient = new Audio("/audio/Cave Wind Loop.mp3");
    ambient.loop = true;
    ambient.volume = 0.3;
    ambientRef.current = ambient;

    return () => {
      music.pause();
      ambient.pause();
    };
  }, []);

  // Switch tracks.
  useEffect(() => {
    const music = musicRef.current;
    const track = TRACKS.find((t) => t.id === currentTrack);
    if (!music || !track) return;
    music.src = track.src;
    music.muted = isMusicMuted;
    music.play().catch(() => {});
  }, [currentTrack]);

  // Keep music mute state in sync.
  useEffect(() => {
    if (musicRef.current) musicRef.current.muted = isMusicMuted;
  }, [isMusicMuted]);

  // Play/pause ambient wind with the effects toggle.
  useEffect(() => {
    const ambient = ambientRef.current;
    if (!ambient) return;
    if (isEffectsMuted) {
      ambient.pause();
    } else {
      ambient.play().catch(() => {});
    }
  }, [isEffectsMuted]);

  const playEffect = (src) => {
    if (isEffectsMuted) return;
    new Audio(src).play().catch(() => {});
  };

  return (
    <AudioCtx.Provider
      value={{
        currentTrack,
        setCurrentTrack,
        isMusicMuted,
        toggleMusicMuted: () => setIsMusicMuted((m) => !m),
        isEffectsMuted,
        toggleEffectsMuted: () => setIsEffectsMuted((m) => !m),
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
