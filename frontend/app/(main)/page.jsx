"use client";

import Link from "next/link";
import { useState } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useAudio } from "@/contexts/AudioContext";
import AvatarImage from "@/components/shared/AvatarImage";

export default function HomePage() {
  return (
    <section className="postsContainer">
      <div className="postsPlaceholder">Your feed will appear here.</div>
    </section>
  );
}
