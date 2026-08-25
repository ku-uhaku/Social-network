import ErrorDisplay from "@/components/ErrorDisplay";



export default function ErrorPage({sherchParams}) {
    const message=sherchParams.message || "something wrong"
  return (
    <ErrorDisplay
      message={message}
      homeHref="/"
    />
  );
}