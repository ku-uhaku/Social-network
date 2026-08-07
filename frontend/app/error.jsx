"use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";
import { useAudio } from "@/contexts/AudioContext";
import { useEffect } from "react";

export default function Error({ error, reset }) {
  const {setMusic} = useAudio();

  useEffect(() => {
    setMusic("/audio/error_music.mp3");
  }, [setMusic])

  return <ErrorDisplay message={error?.message} onRetry={reset} />;
}
