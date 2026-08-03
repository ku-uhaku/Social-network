"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";

export default function LoginPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) router.replace("/");
  }, [loading, user, router]);

  return (
    <div className="page">
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
    <form className="card" onSubmit={handleSubmit}>
      <h1 className="title">Log in</h1>

      {error && <div className="error">{error}</div>}

      <div className="field">
        <label className="label" htmlFor="identifier">Email or nickname</label>
        <input
          id="identifier"
          className="input"
          type="text"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
          required
        />
      </div>

      <div className="field">
        <label className="label" htmlFor="password">Password</label>
        <input
          id="password"
          className="input"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>

      <button className="button" type="submit" disabled={submitting}>
        {submitting ? "Logging in..." : "Log in"}
      </button>

      <div className="footer">
        Don&apos;t have an account? <a className="link" href="/register">Register</a>
      </div>
    </form>
  );
}