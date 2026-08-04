"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import Avatar from "@/components/shared/Avatar";
import CharmToggle from "@/components/shared/CharmToggle";

const MAIN_WALLPAPER = "/images/main_wallpaper.gif";
const MAIN_MUSIC = "/audio/main_music.mp3";
const HEADER_SEPARATOR = "/images/header_separator.png";
const FOOTER_SEPARATOR = "/images/footer_separator.png";
const ICON_SETTINGS = "/images/icon_settings.png";

// TODO: broke settings button
export default function MainLayout({ children }) {
  const { user, loading, logout } = useAuth();
  const { setMusic, isMusicMuted, toggleMusicMuted } = useAudio();
  const router = useRouter();
  const [settingsOpen, setSettingsOpen] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  useEffect(() => {
    setMusic(MAIN_MUSIC);
  }, [setMusic]);

  if (loading || !user) return null;

  const username = user?.username || "User";

  return (
    <main className="homePage">
      <img className="mainWallpaper" src={MAIN_WALLPAPER} alt="" />

      <header className="homeHeader">
        <Link href="/" className="brandLogo">
          Social Network
        </Link>

        <div className="userInfo">
          <Avatar avatar={user?.avatar} name={username} size={52} />
          <div>
            <strong className="userName">{username}</strong>
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

      <img className="headerSeparator" src={HEADER_SEPARATOR} alt="" />

      {/* TODO: settings for particles */}
      <section className={`settingsPanel ${settingsOpen ? "open" : ""}`}>
        <div className="settingsRow">
          <span className="toggleLabel">Music</span>
          <CharmToggle
            checked={!isMusicMuted}
            onChange={() => toggleMusicMuted()}
            label="Music"
          />
        </div>
      </section>

      <section className="pageContent">{children}</section>

      <img className="footerSeparator" src={FOOTER_SEPARATOR} alt="" />

      <footer className="homeFooter">
        <span>Social Network by Mbelhouss and Mbarrah</span>
        <span>Git gud!</span>
      </footer>
    </main>
  );
}