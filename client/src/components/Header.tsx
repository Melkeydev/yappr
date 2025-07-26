import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";

export default function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { showToast } = useToast();

  async function handleLogout() {
    try {
      await logout();
      showToast("Logged out successfully", "success");
      navigate("/login", { replace: true });
    } catch (error: any) {
      showToast(
        "Failed to logout properly. Please clear your browser data if issues persist.",
        "warning",
      );
      navigate("/login", { replace: true });
    }
  }

  return (
    <header className="h-14 flex items-center justify-between px-4 bg-white shadow">
      <div className="flex items-center gap-6">
        <Link
          to="/rooms"
          className="flex items-center hover:opacity-80 transition-opacity"
        >
          <img
            src="/favicon-32x32.png"
            alt="Chat App Logo"
            className="h-8 w-8"
          />
        </Link>
        <Link
          to="/about"
          className="text-sm text-gray-600 hover:text-gray-900 transition-colors"
        >
          About
        </Link>
      </div>

      {user && (
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-600">
            {user.username}
            {user.guest && " (Guest)"}
          </span>
          {!user.guest && (
            <Link
              to="/profile"
              className="text-sm text-indigo-600 hover:text-indigo-500"
            >
              Profile
            </Link>
          )}
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
