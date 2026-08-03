"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import * as authApi from "@/lib/api/auth";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const parseUserResponse = (response) => response?.data ?? response?.user ?? null;

  const refresh = useCallback(async () => {
    try {
      const data = await authApi.me();
      setUser(parseUserResponse(data));
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const login = useCallback(async (identifier, password) => {
    const data = await authApi.login({ identifier, password });
    setUser(parseUserResponse(data));
    return data;
  }, []);

  const register = useCallback(async (formData) => {
    const data = await authApi.register(formData);
    setUser(parseUserResponse(data));
    return data;
  }, []);

  const logout = useCallback(async () => {
    await authApi.logout();
    setUser(null);
  }, []);

  // TODO: simplify `value` later, we don't need all these vars
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