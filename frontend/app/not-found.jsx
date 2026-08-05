
"use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";
import { useAudio } from "@/contexts/AudioContext";
import { useEffect } from "react";

// "not-found" auto-overrides default 404 page
export default function NotFound() {
  const { setMusic } = useAudio();

  useEffect(() => {
    setMusic("/audio/error_music.mp3");
  }, [setMusic]);

  return <ErrorDisplay message="Page not found." />;
}
