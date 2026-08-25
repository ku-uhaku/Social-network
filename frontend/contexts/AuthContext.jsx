"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import * as authApi from "@/lib/api/auth";

const AuthContext = createContext(null);

const parseUserResponse = (response) => response?.data ?? response?.user ?? null;

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  

  const refresh = useCallback(async () => {
    try {
      const data = await authApi.me();
      setUser(parseUserResponse(data));
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;
    authApi.me()
      .then((data) => {
        if (!cancelled) setUser(parseUserResponse(data));
      })
      .catch(() => {
        if (!cancelled) setUser(null);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(async (identifier, password) => {
    const data = await authApi.login({ identifier, password });
    setUser(parseUserResponse(data));
    return data;
  }, []);

  const register = useCallback(async (formData) => {
    return authApi.register(formData);
  }, []);

  const logout = useCallback(async () => {
    await authApi.logout();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout, refresh }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
  return ctx;
}
