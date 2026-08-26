"use client";

import Link from "next/link";
import "@/css/error.css";

export default function ErrorDisplay({ message, onRetry,status, homeHref = "/" }) {
  return (
    <div className="page">
      <img className="wallpaper" src="/images/loading.gif" alt="" />
      <img className="separator" src="/images/card_separator.png" alt="" />
      <div className="card">
        <div className="card-content">
          <h1 className="title">Death and decay awaits</h1>
          <p className="message">{status }</p>
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
        <img className="image" src="/images/err_image.png" alt="" />
      </div>
      <img className="separator separator_right" src="/images/card_separator.png" alt="" />
    </div>
  );
}
