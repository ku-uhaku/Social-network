"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import ImageUploadButton from "@/components/shared/ImageUploadButton";
import AuthBackground from "@/components/shared/AuthBackground";
import NailButton from "@/components/shared/NailButton";
import { useToast } from "@/contexts/ToastContext";
import { isOldEnough } from "@/lib/utils";

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
    <AuthBackground>
      <RegisterForm />
    </AuthBackground>
  );
}

function RegisterForm() {
  const router = useRouter();
  const { register } = useAuth();
  const toassst=useToast()

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
    if (!isOldEnough(values.date_of_birth)){
      setError("You must select a valid age (you must be at least 16 years old)");
      setSubmitting(false);
      return
    }
    const requiredFilled =
      values.email.trim() !== "" &&
      values.first_name.trim() !== "" &&
      values.last_name.trim() !== "" &&
      values.gender.trim() !== "" &&
      values.password.trim() !== "";
    if (!requiredFilled) {
      setError("Required fields cannot be empty or contain only spaces");
      setSubmitting(false);
      return;
    }
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
                  toassst.success("You registereed succesfully now try to login ")
      router.push("/login");
    } catch (err) {
        if (
          err.status >= 500
        ) {
          router.push(
            `/error?message=${(err.statusText)}`
          );
        }
      // setError(err?.message || "Registration failed");
                        toassst.error(err?.message||"You registereed succesfully now try to login")
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="container">
      <img className="separator" src="/images/card_separator.png" alt="" />
      <form className="card" onSubmit={handleSubmit}>
        <h1 className="title">Create an account</h1>
        <div className="field">
          <label className="label" htmlFor="username">Nickname (optional)</label>
          <input id="username" name="username" className="input" type="text" value={values.username} onChange={handleChange} minLength={3} maxLength={20} />
        </div>

        <div className="field">
          <label className="label" htmlFor="email">Email</label>
          <input id="email" name="email" className="input" type="email" value={values.email} onChange={handleChange} required maxLength={254} />
        </div>

        <div className="field">
          <label className="label" htmlFor="password">Password</label>
          <input id="password" name="password" className="input" type="password" value={values.password} onChange={handleChange} required minLength={8} maxLength={72} />
        </div>

        <div className="row">
          <div className="field">
            <label className="label" htmlFor="first_name">First name</label>
            <input id="first_name" name="first_name" className="input" type="text" value={values.first_name} onChange={handleChange} required maxLength={50} />
          </div>
          <div className="field">
            <label className="label" htmlFor="last_name">Last name</label>
            <input id="last_name" name="last_name" className="input" type="text" value={values.last_name} onChange={handleChange} required maxLength={50} />
          </div>
        </div>

        <div className="field">
          <label className="label" htmlFor="gender">Gender</label>
          <select id="gender" name="gender" className="input" value={values.gender} onChange={handleChange} required>
            <option value="">Select gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
          </select>
        </div>

        <div className="field">
          <label className="label" htmlFor="date_of_birth">Date of birth</label>
          <input id="date_of_birth" name="date_of_birth" className="input" type="date" value={values.date_of_birth} onChange={handleChange} required />
        </div>
        <ImageUploadButton
          label="Avatar (optional)"
          value={avatar}
          onChange={setAvatar}
          />
        <div className="field">
          <label className="label" htmlFor="about_me">About me (optional)</label>
          <textarea id="about_me" name="about_me" className="input" rows={3} value={values.about_me} onChange={handleChange} maxLength={500} />
        </div>
          {error && <div className="error">{error}</div>}

        <NailButton type="submit" disabled={submitting}>
          {submitting ? "Creating account..." : "Create account"}
        </NailButton>

        <div className="footer">
          Already have an account? <Link className="link" href="/login">Log in</Link>
        </div>
      </form>
      <img className="separator separator_right" src="/images/card_separator.png" alt="" />
    </div>
  );
}