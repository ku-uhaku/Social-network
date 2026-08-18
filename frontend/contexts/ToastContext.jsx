"use client";

import { createContext, useCallback, useContext, useRef, useState } from "react";
import Link from "next/link";
import "@/css/toast.css";

const ToastContext = createContext(null);

const exxxit_for_animation = 100;

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([]);
  const idRef = useRef(0);
  // we inisial a timer to store tost iddds 
  const timmmerRef = useRef(new Map());

  const removeToast = useCallback((id) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
    const timer = timmmerRef.current.get(id);
    if (timer) {
      clearTimeout(timer);
      timmmerRef.current.delete(id);
    }
  }, []);

  const dismiss = useCallback(
    (id) => {
      const timer = timmmerRef.current.get(id);
      if (timer) {
        clearTimeout(timer);
        timmmerRef.current.delete(id);
      }
      // Mark for exit animation, then remove it.
      setToasts((prev) => prev.map((t) => (t.id === id ? { ...t, leaving: true } : t)));
      setTimeout(() => removeToast(id), exxxit_for_animation);
    },
    [removeToast]
  );

  const dismissAll = useCallback(() => {
    timmmerRef.current.forEach((timer) => clearTimeout(timer));
    timmmerRef.current.clear();
    setToasts((prev) => prev.map((t) => ({ ...t, leaving: true })));
    setTimeout(() => setToasts([]), exxxit_for_animation);
  }, []);

  const show = useCallback(
    (message, options = {}) => {
      const { type = "info", action = null, duration = 5000 } = options;
      const id = ++idRef.current;
      setToasts((prev) => [...prev, { id, type, message, action }]);

      if (duration > 0) {
        const timer = setTimeout(() => dismiss(id), duration);
        timmmerRef.current.set(id, timer);
      }
      return id;
    },
    [dismiss]
  );

  const error = useCallback((message, options) => show(message, { ...options, type: "error" }), [show]);
  const success = useCallback((message, options) => show(message, { ...options, type: "success" }), [show]);
  const info = useCallback((message, options) => show(message, { ...options, type: "info" }), [show]);
  const warning = useCallback((message, options) => show(message, { ...options, type: "warning" }), [show]);

  const value = { show, error, success, info, warning, dismiss, dismissAll };

  return (
    <ToastContext.Provider value={value}>
      {children}

      <div className="toastViewport">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`toast toast-${t.type}${t.leaving ? " toast-leaving" : ""}`}
            role="alert"
          >
            <div className="toastBody">
              <span className="toastMessage">{t.message}</span>
              {t.action && (
                <div className="toastActions">
                  {t.action.href ? (
                    <Link className="toastAction" href={t.action.href} onClick={() => dismiss(t.id)}>
                      {t.action.label}
                    </Link>
                  ) : (
                    <button
                      type="button"
                      className="toastAction"
                      onClick={() => {
                        t.action.onClick?.();
                        dismiss(t.id);
                      }}
                    >
                      {t.action.label}
                    </button>
                  )}
                </div>
              )}
            </div>
            <button
              type="button"
              className="toastClose"
              onClick={() => dismiss(t.id)}
              aria-label="Dismiss"
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within a ToastProvider");
  return ctx;
}
