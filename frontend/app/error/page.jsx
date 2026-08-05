"use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";
import { useAudio } from "@/contexts/AudioContext";
import { useEffect } from "react";

export default function ErrorPage({ searchParams }) {
  const message = searchParams?.message || "Error.";

  const { setMusic } = useAudio();

  useEffect(() => {
    setMusic("/audio/error_music.mp3");
  }, [setMusic]);

  return <ErrorDisplay message={message} />;
}