"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import AvatarImage from "@/components/shared/AvatarImage";

export default function MainLayout({ children }) {
  const { user, loading, logout } = useAuth();
  const { isMusicMuted, toggleMusicMuted, isEffectsMuted, toggleEffectsMuted } = useAudio();
  const router = useRouter();
  const [settingsOpen, setSettingsOpen] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  if (loading || !user) return null;

  const username = user?.username || "User";

  return (
    <main className="homePage">
      <header className="homeHeader">
        <Link href="/" className="brandLogo">
          Social Network
        </Link>

        <div className="userInfo">
          <AvatarImage avatar={user?.avatar} name={username} size={52} />
          <div>
            <strong className="userName">{username}</strong>
          </div>
        </div>

        <div className="headerControls">
          <button
            type="button"
            className="settingsButton"
            onClick={() => setSettingsOpen(!settingsOpen)}
          >
            {settingsOpen ? "Close settings" : "Settings"}
          </button>
          <button type="button" className="logoutButton" onClick={logout}>
            Logout
          </button>
        </div>
      </header>

      <section className={`settingsPanel ${settingsOpen ? "open" : ""}`}>
        <div className="settingsRow">
          <span className="toggleLabel">Music</span>
          <button className="toggleButton" type="button" onClick={toggleMusicMuted}>
            {isMusicMuted ? "Off" : "On"}
          </button>
        </div>

        <div className="settingsRow">
          <span className="toggleLabel">SFX</span>
          <button className="toggleButton" type="button" onClick={toggleEffectsMuted}>
            {isEffectsMuted ? "Off" : "On"}
          </button>
        </div>
      </section>

      <section className="pageContent">{children}</section>

      <footer className="homeFooter">
        <span>Social Network by Mbelhouss and Mbarrah</span>
        <span>Git gud!</span>
      </footer>
    </main>
  );
}
