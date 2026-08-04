"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import { useParticles } from "@/contexts/ParticlesContext";
import Avatar from "@/components/shared/Avatar";
import CharmToggle from "@/components/shared/CharmToggle";

const MAIN_WALLPAPER = "/images/main_wallpaper.gif";
const MAIN_MUSIC = "/audio/main_music.mp3";
const SEPARATOR = "/images/header_separator.png";
const ICON_SETTINGS = "/images/icon_settings.png";

export default function MainLayout({ children }) {
  const { user, loading, logout } = useAuth();
  const router = useRouter();
  const [settingsOpen, setSettingsOpen] = useState(false);
  const headerRef = useRef(null);

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  // Keep height in sync so the settings
  // panel can anchor itself
  useEffect(() => {
    const header = headerRef.current;
    if (!header) return;

    const updateHeight = () => {
      document.documentElement.style.setProperty("--header-height", `${header.offsetHeight}px`);
    };
    updateHeight();

    const observer = new ResizeObserver(updateHeight);
    observer.observe(header);
    return () => observer.disconnect();
  }, []);

  if (loading || !user) return null;

  return (
    <main className="homePage">
      <img className="mainWallpaper" src={MAIN_WALLPAPER} alt="" />

      <header ref={headerRef} className="homeHeader">
        <Link href="/" className="brandLogo">
          Social Network
        </Link>

        <div className="userInfo">
          <Avatar avatar={user?.avatar} username={user.username} size={52} />
          <div>
            <strong className="userName">{user.username}</strong>
          </div>
        </div>

        <div className="headerControls">
          <button
            type="button"
            className="settingsButton"
            title="Settings"
            onClick={() => setSettingsOpen(!settingsOpen)}
          >
            <img className="settingsIcon" src={ICON_SETTINGS} alt="Settings" />
          </button>
          <button type="button" className="logoutButton" onClick={logout}>
            Logout
          </button>
        </div>
      </header>

      <img className="headerSeparator" src={SEPARATOR} alt="" />

      <Settings open={settingsOpen} />

      <section className="pageContent">{children}</section>

      <img className="footerSeparator" src={SEPARATOR} alt="" />

      <footer className="homeFooter">
        <span>Social Network by Mbelhouss and Mbarrah</span>
        <span>Git gud!</span>
      </footer>
    </main>
  );
}

function Settings({ open }) {
  const { setMusic, isMusicMuted, toggleMusicMuted, isSfxMuted, toggleSfxMuted } = useAudio();
  const { isParticlesEnabled, toggleParticles } = useParticles();

  useEffect(() => {
    setMusic(MAIN_MUSIC);
  }, [setMusic]);

  return (
    <section className={`settingsPanel ${open ? "open" : ""}`}>
      <div className="settingsRow">
        <span className="toggleLabel">Music</span>
        <CharmToggle checked={!isMusicMuted} onChange={toggleMusicMuted} title="Music" />
      </div>
      <div className="settingsRow">
        <span className="toggleLabel">Sound</span>
        <CharmToggle checked={!isSfxMuted} onChange={toggleSfxMuted} title="Sound" />
      </div>
      <div className="settingsRow">
        <span className="toggleLabel">Particles</span>
        <CharmToggle checked={isParticlesEnabled} onChange={toggleParticles} title="Particles" />
      </div>
    </section>
  );
}