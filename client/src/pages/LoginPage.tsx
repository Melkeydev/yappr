import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { login } from "../api/auth";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import { useState } from "react";

type Inputs = { email: string; password: string };

export default function LoginPage() {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<Inputs>();

  const { setUser } = useAuth();
  const { showToast } = useToast();
  const [submitError, setSubmitError] = useState<string>("");

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
    setSubmitError("");
    try {
      console.log("Attempting login with:", values.email);
      const user = await login(values.email, values.password);

      setUser(user);
      localStorage.setItem("chat_user", JSON.stringify(user));
      showToast(`Welcome back, ${user.username}!`, "success");

      navigate("/rooms", { replace: true });
    } catch (error: any) {
      console.error("Login error:", error);
      
      let errorMessage = "Login failed";
      
      if (error.response) {
        if (error.response.status === 401) {
          errorMessage = "Invalid email or password";
        } else if (error.response.status >= 500) {
          errorMessage = "Server error. Please try again later.";
        } else {
          errorMessage = error.response.data?.error || `Login failed: ${error.response.status}`;
        }
      } else if (error.request) {
        errorMessage = "Cannot reach server. Please check your connection.";
      } else {
        errorMessage = error.message || "An unexpected error occurred";
      }
      
      setSubmitError(errorMessage);
      showToast(errorMessage, "error");
    }
  }
  /** ──────── UI ──────── */

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-sky-100 to-teal-50 p-4">
      <div className="w-full max-w-4xl flex flex-col lg:flex-row items-center gap-8">
        {/* App Description Section */}
        <div className="flex-1 text-center lg:text-left">
          <div className="mb-8">
            <h1 className="text-4xl lg:text-5xl font-bold text-gray-900 mb-4">
              Anonymous Chat
            </h1>
            <p className="text-xl text-gray-700 mb-6">
              Connect with others in daily topic discussions
            </p>
          </div>
          
          <div className="space-y-4 text-gray-600">
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center mt-0.5">
                <svg className="w-4 h-4 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 13V5a2 2 0 00-2-2H4a2 2 0 00-2 2v8a2 2 0 002 2h3l3 3 3-3h3a2 2 0 002-2zM5 7a1 1 0 011-1h8a1 1 0 110 2H6a1 1 0 01-1-1zm1 3a1 1 0 100 2h3a1 1 0 100-2H6z" clipRule="evenodd" />
                </svg>
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">Daily Fresh Topics</h3>
                <p className="text-sm">New conversation starters every 24 hours from trending tech, news, and interesting facts</p>
              </div>
            </div>
            
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center mt-0.5">
                <svg className="w-4 h-4 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clipRule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">Stay Anonymous</h3>
                <p className="text-sm">Chat freely without revealing your identity. Create an account or join as a guest</p>
              </div>
            </div>
            
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center mt-0.5">
                <svg className="w-4 h-4 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
                </svg>
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">24-Hour Rooms</h3>
                <p className="text-sm">All conversations are temporary - rooms reset daily with fresh topics and clean slates</p>
              </div>
            </div>
          </div>
        </div>

        {/* Login Form Section */}
        <div className="w-full max-w-sm rounded-2xl bg-white/80 shadow-xl backdrop-blur p-8">
          <h2 className="text-2xl font-bold text-center mb-6">Join the Conversation</h2>

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
          {submitError && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md text-sm">
              {submitError}
            </div>
          )}
          
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
    </div>
  );
}
