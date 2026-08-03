"use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";

// "not-found" auto-overrides throwing page
export default function Error({ error, reset }) {
  return <ErrorDisplay message={error?.message} onRetry={reset} />;
}