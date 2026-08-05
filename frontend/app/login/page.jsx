"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import AuthBackground from "@/components/shared/AuthBackground";
import NailButton from "@/components/shared/NailButton";

export default function LoginPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) router.replace("/");
  }, [loading, user, router]);

  return (
    <AuthBackground>
      <LoginForm />
    </AuthBackground>
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
    <div className="container">
      <img className="separator" src="/images/card_separator.png" alt="" />
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
            minLength={8}
          />
        </div>

        <NailButton type="submit" disabled={submitting}>
          {submitting ? "Logging in..." : "Log in"}
        </NailButton>

        <div className="footer">
          {"Don't have an account? "}
          <a className="link" href="/register">Register</a>
        </div>
      </form>
      <img className="separator separator_right" src="/images/card_separator.png" alt="" />
    </div>
  );
}