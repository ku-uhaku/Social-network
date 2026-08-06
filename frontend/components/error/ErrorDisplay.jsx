"use client";

import Link from "next/link";
import "@/css/error.css";

export default function ErrorDisplay({ message, onRetry, homeHref = "/" }) {
  return (
    <div className="page">
      <img className="wallpaper" src="/images/err_wallpaper.png" alt="" />
      <img className="separator" src="/images/card_separator.png" alt="" />
      <div className="card">
        <h1 className="title">Death and decay awaits</h1>
        <p className="message">{message || "An unexpected error occurred."}</p>

        <div className="actions">
          {onRetry && (
            <button className="button" onClick={onRetry}>
              Try again
            </button>
          )}
          <Link className="link" href={homeHref}>Go Home</Link>
        </div>
      </div>
       <img className="separator separator_right" src="/images/card_separator.png" alt="" />
    </div>
  );
}