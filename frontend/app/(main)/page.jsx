"use client";

import Link from "next/link";
import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import AvatarImage from "@/components/shared/AvatarImage";

export default function HomePage() {
  const { user, logout } = useAuth();
  const { isMusicMuted, toggleMusicMuted, isEffectsMuted, toggleEffectsMuted } = useAudio();
  const [settingsOpen, setSettingsOpen] = useState(false);

  const username = user.first_name || user.username || "User";

  return (
    <main className="homePage">
      <header className="homeHeader">
        <Link href="/" className="brandLogo">
          Social Network
        </Link>

        <div className="userInfo">
          <AvatarImage avatar={user?.avatar} name={username} size={52} />
          <div>
            <span className="greeting">Welcome back,</span>
            <strong className="userName">{username}</strong>
          </div>
        </div>

        <div className="headerControls">
          <button type="button" className="logoutButton" onClick={logout}>
            Logout
          </button>
          <button
            type="button"
            className="settingsButton"
            onClick={() => setSettingsOpen(!settingsOpen)}
          >
            {settingsOpen ? "Close settings" : "Settings"}
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

      <section className="postsContainer">
        <div className="postsPlaceholder">Your feed will appear here.</div>
      </section>

      <footer className="homeFooter">
        <span>Social Network by Mbelhouss and Mbarrah</span>
        <span>Git gud!</span>
      </footer>
    </main>
  );
}
