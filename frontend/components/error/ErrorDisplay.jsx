"use client";

import styles from "@/css/error.css";

export default function ErrorDisplay({ message, onRetry, homeHref = "/" }) {
  return (
    <div className={styles.page}>
      <div className={styles.card}>
        <img
          src="/globe.svg"
          alt=""
          className={styles.image}
          width={96}
          height={96}
        />
        <h1 className={styles.title}>Something went wrong</h1>
        <p className={styles.message}>{message || "An unexpected error occurred."}</p>

        <div className={styles.actions}>
          {onRetry && (
            <button className={styles.button} onClick={onRetry}>
              Try again
            </button>
          )}
          <a className={styles.link} href={homeHref}>
            Go home
          </a>
        </div>
      </div>
    </div>
  );
}