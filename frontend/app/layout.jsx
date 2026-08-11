import { AudioProvider } from "@/contexts/AudioContext";
import { AuthProvider } from "@/contexts/AuthContext";
import { ParticlesProvider } from "@/contexts/ParticlesContext";
import { WebSocketProvider } from "@/contexts/WebSocketContext";
import { NotificationProvider } from "@/contexts/NotificationContext";
import SplashScreen from "@/components/shared/SplashScreen";
import "@/css/globals.css";
import "@/css/home.css";
import "@/css/auth.css";
import "@/css/profile.css";
import "@/css/notifications.css";
import "@/css/groups.css";
import "@/css/responsive.css";

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <SplashScreen />
        <AuthProvider>
          <WebSocketProvider>
            <NotificationProvider>
              <AudioProvider>
                <ParticlesProvider>
                  {children}
                </ParticlesProvider>
              </AudioProvider>
            </NotificationProvider>
          </WebSocketProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
