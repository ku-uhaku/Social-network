import { AudioProvider } from "@/contexts/AudioContext";
import { AuthProvider } from "@/contexts/AuthContext";
import { ParticlesProvider } from "@/contexts/ParticlesContext";
import { WebSocketProvider } from "@/contexts/WebSocketContext";
import { NotificationProvider } from "@/contexts/NotificationContext";
import { GroupChatProvider } from "@/contexts/GroupChatContext";
import { ToastProvider } from "@/contexts/ToastContext";
import SplashScreen from "@/components/shared/SplashScreen";
import "@/css/globals.css";
import "@/css/toast.css";
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
        <ToastProvider>
          <AuthProvider>
            <WebSocketProvider>
              <NotificationProvider>
                <GroupChatProvider>
                  <AudioProvider>
                    <ParticlesProvider>
                      {children}
                    </ParticlesProvider>
                  </AudioProvider>
                </GroupChatProvider>
              </NotificationProvider>
            </WebSocketProvider>
          </AuthProvider>
        </ToastProvider>
      </body>
    </html>
  );
}
