"use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";

// TODO: find a hollow knight sprite to put here

// "not-found" auto-overrides throwing page
export default function Error({ error, reset }) {
  return <ErrorDisplay message={error?.message} onRetry={reset} />;
}