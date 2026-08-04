import { AudioProvider } from "@/contexts/AudioContext";
import { AuthProvider } from "@/contexts/AuthContext";
import { ParticlesProvider } from "@/contexts/ParticlesContext";
import "@/css/globals.css";
import "@/css/home.css";
import "@/css/auth.css";

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <AudioProvider>
            <ParticlesProvider>
              {children}
            </ParticlesProvider>
          </AudioProvider>
        </AuthProvider>
      </body>
    </html>
  );
}