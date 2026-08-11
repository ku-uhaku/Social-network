"use client";

import { createContext, useContext, useEffect, useRef, useState, useCallback } from "react";

const CURSOR_IMG = "/images/cursor.png";
const PARTICLE_SPRITES = [
  "/images/mouse_particle_1.png",
  "/images/mouse_particle_2.png",
  "/images/mouse_particle_3.png",
];
const MAX_PARTICLES = 100;
const PARTICLE_LIFETIME = 1000;
const MIN_SIZE = 30;
const MAX_SIZE = 60;

const CRAZY_MAX_PARTICLES = 500;
const CRAZY_MIN_SIZE = 200;
const CRAZY_MAX_SIZE = 500;

const ParticlesCtx = createContext(null);

let particleId = 0;

export function ParticlesProvider({ children }) {
  const [isParticlesEnabled, setIsParticlesEnabled] = useState(true);
  const [isCrazy, setIsCrazy] = useState(false);
  const [particles, setParticles] = useState([]);

  // still need refs for the mousemove closure, but now they just gate spawning
  const settingsRef = useRef({ enabled: isParticlesEnabled, crazy: isCrazy });
  useEffect(() => {
    settingsRef.current = { enabled: isParticlesEnabled, crazy: isCrazy };
  }, [isParticlesEnabled, isCrazy]);

  const removeParticle = useCallback((id) => {
    setParticles((prev) => prev.filter((p) => p.id !== id));
  }, []);

  useEffect(() => {
    document.body.classList.add("customCursor");

    let lastX = 0;
    let lastY = 0;

    function spawnParticle(x, y) {
      const { enabled, crazy } = settingsRef.current;
      if (!enabled) return;

      const size = crazy
        ? CRAZY_MIN_SIZE + Math.random() * (CRAZY_MAX_SIZE - CRAZY_MIN_SIZE)
        : MIN_SIZE + Math.random() * MAX_SIZE;

      const angle = Math.random() * Math.PI * 2;
      const distance = 30 + Math.random() * 40;

      const particle = {
        id: particleId++,
        src: PARTICLE_SPRITES[Math.floor(Math.random() * PARTICLE_SPRITES.length)],
        size,
        x: x - size / 2,
        y: y - size / 2,
        dx: Math.cos(angle) * distance,
        dy: Math.sin(angle) * distance,
      };

      setParticles((prev) => {
        const maxParticles = crazy ? CRAZY_MAX_PARTICLES : MAX_PARTICLES;
        const next = [...prev, particle];
        return next.length > maxParticles ? next.slice(next.length - maxParticles) : next;
      });

      setTimeout(() => removeParticle(particle.id), PARTICLE_LIFETIME);
    }

    function handleMove(e) {
      const cursor = document.getElementById("customCursorImg");
      if (cursor) {
        cursor.style.left = `${e.clientX}px`;
        cursor.style.top = `${e.clientY}px`;
      }

      const dx = e.clientX - lastX;
      const dy = e.clientY - lastY;
      const distance = Math.sqrt(dx * dx + dy * dy);
      if (distance > (settingsRef.current.crazy ? 0.2 : 6)) {
        spawnParticle(e.clientX, e.clientY);
        lastX = e.clientX;
        lastY = e.clientY;
      }
    }

    window.addEventListener("mousemove", handleMove);

    return () => {
      window.removeEventListener("mousemove", handleMove);
      document.body.classList.remove("customCursor");
    };
  }, [removeParticle]);

  return (
    <ParticlesCtx.Provider
      value={{
        isParticlesEnabled,
        toggleParticles: () => setIsParticlesEnabled((e) => !e),
        isCrazy,
        toggleCrazy: () => setIsCrazy((c) => !c),
      }}
    >
      {children}
      <div className="cursorOverlay">
        <img id="customCursorImg" className="cursorImage" src={CURSOR_IMG} alt="" />
        {particles.map((p) => (
          <img
            key={p.id}
            className="cursorParticle"
            src={p.src}
            alt=""
            style={{
              width: p.size,
              height: p.size,
              left: p.x,
              top: p.y,
              "--dx": `${p.dx}px`,
              "--dy": `${p.dy}px`,
            }}
          />
        ))}
      </div>
    </ParticlesCtx.Provider>
  );
}

export function useParticles() {
  const ctx = useContext(ParticlesCtx);
  if (!ctx) throw new Error("useParticles must be used within a ParticlesProvider");
  return ctx;
}