"use client";

import { useEffect } from "react";
import { useAudio } from "@/contexts/AudioContext";

const AUTH_WALLPAPER = "/images/auth_wallpaper.gif";
const AUTH_MUSIC = "/audio/auth_music.mp3";

export default function AuthBackground({ children }) {
  const { setMusic } = useAudio();

  // TODO: why is setMusic always inside useEffect?
  useEffect(() => {
    setMusic(AUTH_MUSIC);
  }, [setMusic]);

  return (
    <div className="authPage">
      <img className="authWallpaper" src={AUTH_WALLPAPER} alt="" />
      <div className="authContent">{children}</div>
    </div>
  );
}