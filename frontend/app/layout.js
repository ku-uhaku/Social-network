import { AudioProvider } from "@/contexts/AudioContext";
import "@/css/globals.css";

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AudioProvider>
          {children}
        </AudioProvider>
      </body>
    </html>
  );
}
