"use client";

import { useId, useRef, useState } from "react";
import "@/css/image-upload.css";

const ACCEPTED_TYPES = ["image/png", "image/jpeg", "image/gif"];
const ACCEPTED_EXT = ".png,.jpg,.jpeg,.gif";
const MAX_SIZE_MB = 20;

export default function ImageUploadButton({ label = "Upload image", value, onChange, required = false }) {
  const inputId = useId();
  const inputRef = useRef(null);
  const [error, setError] = useState("");
  const [preview, setPreview] = useState(null);

  function reset() {
    if (inputRef.current) inputRef.current.value = "";
    setPreview(null);
    onChange?.(null);
  }

  function handleFile(e) {
    const file = e.target.files?.[0];
    setError("");

    if (!file) {
      reset();
      return;
    }

    if (!ACCEPTED_TYPES.includes(file.type)) {
      setError("Only PNG, JPEG, or GIF images are allowed.");
      reset();
      return;
    }

    if (file.size > MAX_SIZE_MB * 1024 * 1024) {
      setError(`Image must be smaller than ${MAX_SIZE_MB}MB.`);
      reset();
      return;
    }

    setPreview(URL.createObjectURL(file));
    onChange?.(file);
  }

  return (
    <div className="field">
      <label className="label" htmlFor={inputId}>{label}</label>

      <div className="row">
        <input
          id={inputId}
          ref={inputRef}
          className="input"
          type="file"
          accept={ACCEPTED_EXT}
          onChange={handleFile}
          required={required}
        />
        {(preview || value) && (
          <button type="button" className="clearButton" onClick={reset}>
            Remove
          </button>
        )}
      </div>

      {preview && (
        <img src={preview} alt="Preview" className="preview" />
      )}

      {error && <div className="error">{error}</div>}
    </div>
  );
}