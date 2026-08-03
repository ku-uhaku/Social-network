"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import styles from "@/css/auth.css";

export default function LoginPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) router.replace("/");
  }, [loading, user, router]);

  return (
    <div className={styles.page}>
      <LoginForm />
    </div>
  );
}

function LoginForm() {
  const router = useRouter();
  const { login } = useAuth();

  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      await login(identifier, password);
      router.push("/");
    } catch (err) {
      setError(err?.message || "Login failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className={styles.card} onSubmit={handleSubmit}>
      <h1 className={styles.title}>Log in</h1>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.field}>
        <label className={styles.label} htmlFor="identifier">Email or nickname</label>
        <input
          id="identifier"
          className={styles.input}
          type="text"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
          required
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="password">Password</label>
        <input
          id="password"
          className={styles.input}
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>

      <button className={styles.button} type="submit" disabled={submitting}>
        {submitting ? "Logging in..." : "Log in"}
      </button>

      <div className={styles.footer}>
        Don&apos;t have an account? <a className={styles.link} href="/register">Register</a>
      </div>
    </form>
  );
}