import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { login } from "../api/auth";
import { useAuth } from "../context/AuthContext";

type Inputs = { email: string; password: string };

export default function LoginPage() {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<Inputs>();

  const { setUser } = useAuth();

  /** ──────── handlers ──────── */

  // guest path
  function handleGuest() {
    const username = `anonymousUser_${Math.random().toString(36).slice(-6)}`;
    const guestUser = { id: username, username, guest: true };

    setUser(guestUser); // NEW
    localStorage.setItem("chat_user", JSON.stringify(guestUser));

    navigate("/rooms", { replace: true });
  }

  // real login

  async function onSubmit(values: Inputs) {
    try {
      const user = await login(values.email, values.password);

      setUser(user); // NEW
      localStorage.setItem("chat_user", JSON.stringify(user));

      navigate("/rooms", { replace: true });
    } catch {
      alert("Login failed");
    }
  }
  /** ──────── UI ──────── */

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-sky-100 to-teal-50">
      <div className="w-full max-w-sm rounded-2xl bg-white/80 shadow-xl backdrop-blur p-8">
        <h1 className="text-2xl font-bold text-center mb-6">Welcome to Chat</h1>

        {/* Guest button */}
        <button
          onClick={handleGuest}
          className="w-full py-2 mb-4 text-sm font-medium bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition"
        >
          Continue as guest
        </button>

        {/* Divider */}
        <div className="relative my-4">
          <span className="absolute inset-0 flex items-center justify-center">
            <span className="h-px w-full bg-gray-300" />
          </span>
          <span className="relative bg-white px-2 text-xs text-gray-500">
            OR
          </span>
        </div>

        {/* Login form */}
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              type="email"
              {...register("email", { 
                required: "Email is required",
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: "Invalid email address"
                }
              })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.email && (
              <p className="mt-1 text-xs text-red-600">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              type="password"
              {...register("password", { 
                required: "Password is required",
                minLength: {
                  value: 6,
                  message: "Password must be at least 6 characters"
                }
              })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.password && (
              <p className="mt-1 text-xs text-red-600">{errors.password.message}</p>
            )}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full py-2 text-sm font-medium bg-gray-800 hover:bg-gray-900 text-white rounded-lg"
          >
            {isSubmitting ? "Logging in…" : "Log in"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Don't have an account?{" "}
          <a href="/signup" className="text-indigo-600 hover:text-indigo-500">
            Sign up
          </a>
        </p>
      </div>
    </div>
  );
}
