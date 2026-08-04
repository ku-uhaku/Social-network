"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";

const AudioCtx = createContext(null);

export function AudioProvider({ children }) {
  const musicRef = useRef(null);
  const musicSrcRef = useRef(null);
  const mutedRef = useRef(false);
  const [musicSrc, setMusicSrc] = useState(null);
  const [isMusicMuted, setIsMusicMuted] = useState(false);

  // keep refs in sync with latest state (refs only touched inside effects)
  useEffect(() => {
    musicSrcRef.current = musicSrc;
    mutedRef.current = isMusicMuted;
  }, [musicSrc, isMusicMuted]);

  useEffect(() => {
    const music = new Audio();
    music.loop = true;
    musicRef.current = music;

    const unlock = () => { // sound if muted by default until poage is clicked
      const el = musicRef.current;
      if (el && musicSrcRef.current) {
        el.muted = mutedRef.current;
        el.play().catch(() => {});
      }
      window.removeEventListener("pointerdown", unlock);
      window.removeEventListener("keydown", unlock);
    };
    window.addEventListener("pointerdown", unlock);
    window.addEventListener("keydown", unlock);

    return () => {
      music.pause();
      window.removeEventListener("pointerdown", unlock);
      window.removeEventListener("keydown", unlock);
    };
  }, []);

  useEffect(() => {
    const music = musicRef.current;
    if (!music || !musicSrc) return;
    music.src = musicSrc;
    music.load();
    music.muted = isMusicMuted;
    music.play().catch(() => {});
  }, [musicSrc, isMusicMuted]);

  useEffect(() => {
    if (musicRef.current) musicRef.current.muted = isMusicMuted;
  }, [isMusicMuted]);

  return (
    <AudioCtx.Provider
      value={{
        setMusic: setMusicSrc,
        isMusicMuted,
        toggleMusicMuted: () => setIsMusicMuted((m) => !m),
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