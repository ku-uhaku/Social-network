import ErrorDisplay from "@/components/error/ErrorDisplay";

export default async function ErrorPage({ searchParams }) {
  const params = await searchParams;
  const message = params?.message || "An unexpected error occurred.";

  return <ErrorDisplay message={message} />;
}