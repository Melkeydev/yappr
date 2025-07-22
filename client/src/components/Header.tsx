import { useNavigate, Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";

export default function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { showToast } = useToast();

  async function handleLogout() {
    try {
      await logout(); // clears cookie + localStorage
      showToast("Logged out successfully", "success");
      navigate("/login", { replace: true });
    } catch (error: any) {
      console.error("Logout error:", error);
      showToast("Failed to logout properly. Please clear your browser data if issues persist.", "warning");
      // Still navigate even if logout failed
      navigate("/login", { replace: true });
    }
  }

  return (
    <header className="h-14 flex items-center justify-between px-4 bg-white shadow">
      <Link to="/rooms" className="text-lg font-semibold hover:text-gray-700">
        Chat App
      </Link>

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
