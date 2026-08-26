"use client";

import Link from 'next/link';
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import AuthBackground from "@/components/shared/AuthBackground";
import NailButton from "@/components/shared/NailButton";
import { useToast } from '@/contexts/ToastContext';

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
//use tossast 
const toooasst=useToast()
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    if (identifier.trim() === "" || password.trim() === "") return;
    setSubmitting(true);
    try {
      await login(identifier, password);
        toooasst.success("You welcooome again to our social netwooook ")
      router.push("/");
    } catch (err) {
       
        console.log(err.status)
      // setError(err?.message || "Login failed");
            toooasst.error(err?.message ||"Login failed")
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
            maxLength={254}
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
            maxLength={72}
          />
        </div>

        <NailButton type="submit" disabled={submitting}>
          {submitting ? "Logging in..." : "Log in"}
        </NailButton>

        <div className="footer">
          {"Don't have an account? "}
          <Link className="link" href="/register">Register</Link>
        </div>
      </form>
      <img className="separator separator_right" src="/images/card_separator.png" alt="" />
    </div>
  );
}