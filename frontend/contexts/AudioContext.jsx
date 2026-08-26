"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";

const AudioCtx = createContext(null);

export function AudioProvider({ children }) {
  const musicRef = useRef(null);
  const mutedRef = useRef(false);
  const router=useRouter()
  const [musicSrc, setMusicSrc] = useState(null);
  const [isMusicMuted, setIsMusicMuted] = useState(false);
  const [isSfxMuted, setIsSfxMuted] = useState(false);

  useEffect(() => {
    mutedRef.current = isMusicMuted;
  }, [isMusicMuted]);

  useEffect(() => {
    const music = new Audio();
    music.loop = true;
    musicRef.current = music;

    // sound is muted by default until page is clicked
    const unlock = () => {
      const el = musicRef.current;
      if (el && el.src) {
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
    if (!music) return;
    music.muted = isMusicMuted;
    if (musicSrc) {
      music.src = musicSrc;
      music.load();
      music.play().catch(() => {});
    }
  }, [musicSrc, isMusicMuted]);

  function playSfx(src) {
    if (isSfxMuted || !src) return;
    const sfx = new Audio(src);
    sfx.play().catch(() => {});
  }

  return (
    <AudioCtx.Provider
      value={{
        setMusic: setMusicSrc,
        isMusicMuted,
        toggleMusicMuted: () => setIsMusicMuted((m) => !m),
        isSfxMuted,
        toggleSfxMuted: () => setIsSfxMuted((m) => !m),
        playSfx,
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