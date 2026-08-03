"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import styles from "@/css/auth.css";

const initialState = {
  username: "",
  email: "",
  password: "",
  first_name: "",
  last_name: "",
  gender: "",
  date_of_birth: "",
  about_me: "",
};

export default function RegisterPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) router.replace("/");
  }, [loading, user, router]);

  return (
    <div className={styles.page}>
      <RegisterForm />
    </div>
  );
}

function RegisterForm() {
  const router = useRouter();
  const { register } = useAuth();

  const [values, setValues] = useState(initialState);
  const [avatar, setAvatar] = useState(null);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  function handleChange(e) {
    const { name, value } = e.target;
    setValues((v) => ({ ...v, [name]: value }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    setSubmitting(true);

    try {
      const formData = new FormData();
      formData.append("username", values.username.trim());
      formData.append("email", values.email.trim());
      formData.append("password", values.password);
      formData.append("first_name", values.first_name.trim());
      formData.append("last_name", values.last_name.trim());
      formData.append("gender", values.gender);
      formData.append("date_of_birth", values.date_of_birth);
      if (values.about_me?.trim()) formData.append("about_me", values.about_me.trim());
      if (avatar) formData.append("avatar", avatar);

      await register(formData);
      router.push("/login");
    } catch (err) {
      setError(err?.message || "Registration failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className={styles.card} onSubmit={handleSubmit}>
      <h1 className={styles.title}>Create an account</h1>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.field}>
        <label className={styles.label} htmlFor="username">Nickname</label>
        <input id="username" name="username" className={styles.input} type="text" value={values.username} onChange={handleChange} required />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="email">Email</label>
        <input id="email" name="email" className={styles.input} type="email" value={values.email} onChange={handleChange} required />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="password">Password</label>
        <input id="password" name="password" className={styles.input} type="password" value={values.password} onChange={handleChange} required minLength={8} />
      </div>

      <div className={styles.row}>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="first_name">First name</label>
          <input id="first_name" name="first_name" className={styles.input} type="text" value={values.first_name} onChange={handleChange} required />
        </div>
        <div className={styles.field}>
          <label className={styles.label} htmlFor="last_name">Last name</label>
          <input id="last_name" name="last_name" className={styles.input} type="text" value={values.last_name} onChange={handleChange} required />
        </div>
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="gender">Gender</label>
        <select id="gender" name="gender" className={styles.input} value={values.gender} onChange={handleChange} required>
          <option value="">Select gender</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
          <option value="other">Other</option>
        </select>
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="date_of_birth">Date of birth</label>
        <input id="date_of_birth" name="date_of_birth" className={styles.input} type="date" value={values.date_of_birth} onChange={handleChange} required />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="avatar">Avatar (optional)</label>
        <input
          id="avatar"
          name="avatar"
          className={styles.input}
          type="file"
          accept="image/jpeg,image/png,image/gif"
          onChange={(e) => setAvatar(e.target.files?.[0] ?? null)}
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label} htmlFor="about_me">About me (optional)</label>
        <textarea id="about_me" name="about_me" className={styles.input} rows={3} value={values.about_me} onChange={handleChange} />
      </div>

      <button className={styles.button} type="submit" disabled={submitting}>
        {submitting ? "Creating account..." : "Create account"}
      </button>

      <div className={styles.footer}>
        Already have an account? <a className={styles.link} href="/login">Log in</a>
      </div>
    </form>
  );
}