import { API_BASE } from "@/lib/utils";

 export  async function apiFetch(path, options = {}) {
   const { body, headers, ...rest } = options;
   const isFormData = body instanceof FormData;
   
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: isFormData ? headers : { "Content-Type": "application/json", ...headers },
    body: isFormData ? body : body ? JSON.stringify(body) : undefined,
    ...rest,
  });
  if (!res.ok){
    if (
          res.status === 401 ||
          res.status === 403 ||
          res.status === 404 ||
          res.status >= 500
        ) {
          throw res
        }
  }
  
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
    const details = Array.isArray(data?.errors)
      ? `: ${data.errors
        .map((e) => (typeof e === "string" ? e : e?.message || JSON.stringify(e)))
        .join(", ")}`
      : "";
    throw new Error(`${data?.message || res.statusText}${details}`);
  }

  return data;
}