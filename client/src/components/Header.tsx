import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    await logout(); // clears cookie + localStorage
    navigate("/login", { replace: true });
  }

  return (
    <header className="h-14 flex items-center justify-between px-4 bg-white shadow">
      <span className="text-lg font-semibold">Chat App</span>

      {user && (
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-600">
            {user.username}
            {user.guest && " (Guest)"}
          </span>
          <button
            onClick={handleLogout}
            className="rounded-md bg-gray-800 px-3 py-1 text-xs font-medium text-white hover:bg-gray-900 transition"
          >
            Logout
          </button>
        </div>
      )}
    </header>
  );
}
