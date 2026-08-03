
import { redirect } from "next/navigation";

// "not-found" auto-overrides default 404 page
export default function NotFound() {
  redirect("/error?message=" + encodeURIComponent("Page not found."));
}