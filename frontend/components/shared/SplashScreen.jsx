"use client";

import { useEffect, useState } from "react";

const SHOW_MS = 1000;
const FADE_MS = 500;

export default function SplashScreen() {
  const [visible, setVisible] = useState(true);
  const [dismissed, setDismissed] = useState(false);

  useEffect(() => {
    const fadeTimer = setTimeout(() => setVisible(false), SHOW_MS);
    const removeTimer = setTimeout(() => setDismissed(true), SHOW_MS + FADE_MS);
    return () => {
      clearTimeout(fadeTimer);
      clearTimeout(removeTimer);
    };
  }, []);

  if (dismissed) return null;

  // transition handled in css
  return (
    <div className={`splashScreen ${visible ? "" : "hidden"}`}>
      <img className="splashLogo" src="/images/main_logo.png" alt="Social Network" />
    </div>
  );
}
