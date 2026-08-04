import { AudioProvider } from "@/contexts/AudioContext";
import { AuthProvider } from "@/contexts/AuthContext";
import CursorParticles from "@/components/shared/CursorParticles";
import "@/css/globals.css";
import "@/css/home.css";
import "@/css/auth.css";

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <CursorParticles />
        <AuthProvider>
          <AudioProvider>
            {children}
          </AudioProvider>
        </AuthProvider>
      </body>
    </html>
  );
}