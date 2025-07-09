import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import LoginPage from "./pages/LoginPage";
import ChatPage from "./pages/ChatPage";
import Protected from "./components/Protected";
import RoomsPage from "./pages/RoomPage";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />

          <Route element={<Protected />}>
            <Route path="/rooms" element={<RoomsPage />} />
            <Route path="/chat/:roomId" element={<ChatPage />} />
            <Route index element={<Navigate to="/rooms" replace />} />
          </Route>

          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
