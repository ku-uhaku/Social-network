"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import { useParticles } from "@/contexts/ParticlesContext";
import { useNotifications } from "@/contexts/NotificationContext";
import Avatar from "@/components/shared/Avatar";
import CharmToggle from "@/components/shared/CharmToggle";
import IconButton from "@/components/shared/IconButton";
import NotificationList from "@/components/notifications/NotificationList";
import Chat, { useChat } from "@/components/chat/Chat";

const MAIN_WALLPAPER = "/images/main_wallpaper.png";
const MAIN_MUSIC = "/audio/main_music.mp3";
const SEPARATOR = "/images/header_separator.png";
const ICON_SETTINGS = "/images/icon_settings.png";
const ICON_NOTIFICATION_ON = "/images/notification_on.png";
const ICON_NOTIFICATION_OFF = "/images/notification_off.png";
const ICON_LOGOUT = "/images/icon_logout.png";

export default function MainLayout({ children }) {
  const { user, loading, logout } = useAuth();
  const { unreadCount } = useNotifications();
  const router = useRouter();
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const headerRef = useRef(null);
  const chat = useChat();

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  // Keep height in sync so the settings
  // panel can anchor itself
  useEffect(() => { // TODO: this whole function stinks
    if (loading || !user) return;

    const header = headerRef.current;
    if (!header) return;

    const updateHeight = () => {
      document.documentElement.style.setProperty("--header-height", `${header.offsetHeight}px`);
    };
    updateHeight();

    const observer = new ResizeObserver(updateHeight);
    observer.observe(header);
    return () => observer.disconnect();
  }, [loading, user]);

  if (loading || !user) return null;

  return (
    <main className="homePage">
      <img className="mainWallpaper" src={MAIN_WALLPAPER} alt="" />

      <header ref={headerRef} className="homeHeader">
        <div className="headerLeft">
          <Link href="/" className="brandLogo">
            <img src="/images/main_logo.png" alt="Social Network" />
          </Link>

          <span className="headerDividerV" />

          <div className="userInfo">
            <Avatar avatar={user?.avatar} username={user.username} size={52} />
            <div>
              <strong className="userName">{user.username}</strong>
            </div>
          </div>
        </div>

        <div className="headerRight">
          <div className="headerRightRow">
            <div className="headerControlWrap">
              <IconButton
                icon={unreadCount > 0 ? ICON_NOTIFICATION_ON : ICON_NOTIFICATION_OFF}
                label="Notifications"
                onClick={() => setNotificationsOpen(!notificationsOpen)}
              >
                {unreadCount > 0 && <span className="iconButtonBadge">{unreadCount}</span>}
              </IconButton>
              <NotificationList open={notificationsOpen} />
            </div>
            <div className="headerControlWrap">
              <IconButton
                icon={ICON_SETTINGS}
                label="Settings"
                onClick={() => setSettingsOpen(!settingsOpen)}
              />
              <Settings open={settingsOpen} />
            </div>
            <IconButton icon={ICON_LOGOUT} label="Logout" onClick={logout} />
          </div>

          <span className="headerDividerH" />

          <div className="headerRightRow">
            <IconButton
              icon="/images/icon_chat.png"
              label="Chat"
              active={chat.open}
              onClick={() => chat.setOpen((open) => !open)}
            >
              {chat.totalUnread > 0 && <span className="iconButtonBadge">{chat.totalUnread}</span>}
            </IconButton>
            <IconButton href="/group" label="Groups" icon="/images/icon_groups.png" />
          </div>
        </div>
      </header>

      <img className="headerSeparator" src={SEPARATOR} alt="" />

      <Chat chat={chat} />

      <section className="pageContent">{children}</section>

      <img className="footerSeparator" src={SEPARATOR} alt="" />

      <footer className="homeFooter">
        <span>Social Network by <span className="credits">Mbelhouss</span> and <span className="credits">Mbarrah</span>. Git gud!</span>
      </footer>
    </main>
  );
}

function Settings({ open }) {
  const { setMusic, isMusicMuted, toggleMusicMuted, isSfxMuted, toggleSfxMuted } = useAudio();
  const { isParticlesEnabled, toggleParticles, isCrazy, toggleCrazy } = useParticles();

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
        <div className="particlesControls">
          <CharmToggle checked={isParticlesEnabled} onChange={toggleParticles} title="Particles" />
          <CharmToggle checked={isCrazy} onChange={toggleCrazy} title="???" />
        </div>
      </div>
    </section>
  );
}
