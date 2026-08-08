"use client";

import Link from "next/link";

function DockButton({ href, title, icon, active, onClick, children }) {
  const className = `chatDockButton ${active ? "active" : ""}`;
  const content = (
    <>
      <img className="chatDockIcon" src={icon} alt={title} />
      {children}
    </>
  );

  if (href) {
    return (
      <Link href={href} className={className} title={title}>
        {content}
      </Link>
    );
  }

  return (
    <button type="button" className={className} title={title} onClick={onClick}>
      {content}
    </button>
  );
}

export default function SideBar({ chat }) {
  return (
    <aside className="chatDock">
      <DockButton
        title="Chat"
        icon="/images/icon_chat.png"
        active={chat?.open}
        onClick={() => chat?.setOpen((open) => !open)}
      >
        {chat?.totalUnread > 0 && <span className="chatDockBadge">{chat.totalUnread}</span>}
      </DockButton>

      <DockButton href="/group" title="Groups" icon="/images/icon_chat.png" />
    </aside>
  );
}
