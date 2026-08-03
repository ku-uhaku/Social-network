const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function apiFetch(path, options = {}) {
  const { body, headers, ...rest } = options;
  const isFormData = body instanceof FormData;

  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: isFormData ? headers : { "Content-Type": "application/json", ...headers },
    body: isFormData ? body : body ? JSON.stringify(body) : undefined,
    ...rest,
  });

  const text = await res.text();
  let data = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
  }

  if (!res.ok) {
    throw new Error((data && data.message) || res.statusText);
  }

  return data;
}