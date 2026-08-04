"use client";

import { useEffect, useRef } from "react";

const CURSOR_IMG = "/images/cursor.png";
const PARTICLES = [
  "/images/mouse_particle_1.png",
  "/images/mouse_particle_2.png",
  "/images/mouse_particle_3.png",
];
const MAX_PARTICLES = 24;
const PARTICLE_LIFETIME = 700;

export default function CursorParticles() {
  const overlayRef = useRef(null);
  const cursorRef = useRef(null);

  useEffect(() => {
    document.body.classList.add("customCursor");
    const overlay = overlayRef.current;
    const cursor = cursorRef.current;
    if (!overlay || !cursor) return;

    let particles = [];
    let lastX = 0;
    let lastY = 0;

    function spawnParticle(x, y) {
      if (particles.length >= MAX_PARTICLES) {
        const oldest = particles.shift();
        oldest.remove();
      }

      const img = document.createElement("img");
      img.className = "cursorParticle";
      img.src = PARTICLES[Math.floor(Math.random() * PARTICLES.length)];
      const size = 14 + Math.random() * 22;
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
      if (distance > 12) {
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
    <div className="cursorOverlay" ref={overlayRef}>
      <img className="cursorImage" ref={cursorRef} src={CURSOR_IMG} alt="" />
    </div>
  );
}