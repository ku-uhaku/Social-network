// "use client";

import ErrorDisplay from "@/components/error/ErrorDisplay";
// import { useAudio } from "@/contexts/AudioContext";
// import { useEffect } from "react";

export default async function ErrorPage({ searchParams }) {
  const paramas = await searchParams ;
  console.log("msg::",paramas);
  
  // const { setMusic } = useAudio();

  // useEffect(() => {
  //   setMusic("/audio/error_music.mp3");
  // }, [setMusic]);

  return <ErrorDisplay message={paramas.message} status={paramas.status} />;
}
