"use client";

import { createContext, useContext, useEffect, useRef, useState } from "react";

const CURSOR_IMG = "/images/cursor.png";
const PARTICLE_SPRITES = [
  "/images/mouse_particle_1.png",
  "/images/mouse_particle_2.png",
  "/images/mouse_particle_3.png",
];
const MAX_PARTICLES = 100;
const PARTICLE_LIFETIME = 1000;

const ParticlesCtx = createContext(null);

export function ParticlesProvider({ children }) {
  const [isParticlesEnabled, setIsParticlesEnabled] = useState(true);
  const enabledRef = useRef(isParticlesEnabled);

  const overlayRef = useRef(null);
  const cursorRef = useRef(null);

  useEffect(() => {
    enabledRef.current = isParticlesEnabled;
  }, [isParticlesEnabled]);

  useEffect(() => {
    document.body.classList.add("customCursor");
    const overlay = overlayRef.current;
    const cursor = cursorRef.current;
    if (!overlay || !cursor) return;

    let particles = [];
    let lastX = 0;
    let lastY = 0;

    function spawnParticle(x, y) {
      if (!enabledRef.current) return;

      if (particles.length >= MAX_PARTICLES) {
        const oldest = particles.shift();
        oldest.remove();
      }

      const img = document.createElement("img");
      img.className = "cursorParticle";
      img.src = PARTICLE_SPRITES[Math.floor(Math.random() * PARTICLE_SPRITES.length)];
      const size = 30 + Math.random() * 60;
      img.style.width = `${size}px`;
      img.style.height = `${size}px`;

      const angle = Math.random() * Math.PI * 2;
      const distance = 30 + Math.random() * 40;
      img.style.setProperty("--dx", `${Math.cos(angle) * distance}px`);
      img.style.setProperty("--dy", `${Math.sin(angle) * distance}px`);

      img.style.left = `${x - size / 2}px`;
      img.style.top = `${y - size / 2}px`;
      overlay.appendChild(img);

      particles.push(img);
      setTimeout(() => img.remove(), PARTICLE_LIFETIME);
    }

    function handleMove(e) {
      cursor.style.left = `${e.clientX}px`;
      cursor.style.top = `${e.clientY}px`;

      const dx = e.clientX - lastX;
      const dy = e.clientY - lastY;
      const distance = Math.sqrt(dx * dx + dy * dy);
      if (distance > 6) {
        spawnParticle(e.clientX, e.clientY);
        lastX = e.clientX;
        lastY = e.clientY;
      }
    }

    window.addEventListener("mousemove", handleMove);

    return () => {
      window.removeEventListener("mousemove", handleMove);
      document.body.classList.remove("customCursor");
      particles.forEach((p) => p.remove());
      particles = [];
    };
  }, []);

  return (
    <ParticlesCtx.Provider
      value={{
        isParticlesEnabled,
        toggleParticles: () => setIsParticlesEnabled((e) => !e),
      }}
    >
      {children}
      <div className="cursorOverlay" ref={overlayRef}>
        <img className="cursorImage" ref={cursorRef} src={CURSOR_IMG} alt="" />
      </div>
    </ParticlesCtx.Provider>
  );
}

export function useParticles() {
  const ctx = useContext(ParticlesCtx);
  if (!ctx) throw new Error("useParticles must be used within a ParticlesProvider");
  return ctx;
}