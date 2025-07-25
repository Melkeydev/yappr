import { useForm } from "react-hook-form";
import { useNavigate, Link } from "react-router-dom";
import { signup } from "../api/auth";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import { useState, useEffect } from "react";

type Inputs = {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
};

export default function SignupPage() {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<Inputs>();

  const { user, setUser } = useAuth();
  const { showToast } = useToast();
  const [submitError, setSubmitError] = useState<string>("");

  // Redirect if already logged in
  useEffect(() => {
    if (user) {
      navigate("/rooms", { replace: true });
    }
  }, [user, navigate]);

  async function onSubmit(values: Inputs) {
    setSubmitError("");
    try {
      const user = await signup(values.username, values.email, values.password);
      setUser(user);
      localStorage.setItem("chat_user", JSON.stringify(user));
      showToast(
        `Welcome, ${user.username}! Account created successfully.`,
        "success",
      );
      navigate("/rooms", { replace: true });
    } catch (error: any) {
      let errorMessage = "Signup failed";

      if (error.response) {
        // Server responded with error status
        if (
          error.response.status === 400 &&
          error.response.data?.error?.includes("inappropriate content")
        ) {
          errorMessage =
            "Username contains inappropriate content. Please choose a different username.";
        } else {
          errorMessage =
            error.response.data?.error ||
            `Server error: ${error.response.status}`;
        }
      } else if (error.request) {
        // Request was made but no response received (network issues, CORS, etc.)
        errorMessage = "Cannot reach server. Please check your connection.";
      } else {
        // Something else happened
        errorMessage = error.message || "An unexpected error occurred";
      }

      setSubmitError(errorMessage);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-sky-100 to-teal-50">
      <div className="w-full max-w-sm rounded-2xl bg-white/80 shadow-xl backdrop-blur p-8">
        <h1 className="text-2xl font-bold text-center mb-6">Create Account</h1>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {submitError && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md text-sm">
              {submitError}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Username
            </label>
            <input
              type="text"
              {...register("username", {
                required: "Username is required",
                minLength: {
                  value: 3,
                  message: "Username must be at least 3 characters",
                },
                maxLength: {
                  value: 20,
                  message: "Username must be less than 20 characters",
                },
                pattern: {
                  value: /^[a-zA-Z0-9_]+$/,
                  message:
                    "Username can only contain letters, numbers, and underscores",
                },
              })}
              className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.username && (
              <p className="mt-1 text-xs text-red-600">
                {errors.username.message}
              </p>
            )}
          </div>

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
                  message: "Invalid email address",
                },
              })}
              className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.email && (
              <p className="mt-1 text-xs text-red-600">
                {errors.email.message}
              </p>
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
                  message: "Password must be at least 6 characters",
                },
              })}
              className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.password && (
              <p className="mt-1 text-xs text-red-600">
                {errors.password.message}
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Confirm Password
            </label>
            <input
              type="password"
              {...register("confirmPassword", {
                required: "Please confirm your password",
                validate: (value) =>
                  value === watch("password") || "Passwords don't match",
              })}
              className="mt-1 block w-full rounded-md border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
            {errors.confirmPassword && (
              <p className="mt-1 text-xs text-red-600">
                {errors.confirmPassword.message}
              </p>
            )}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full py-2 text-sm font-medium bg-gray-800 hover:bg-gray-900 text-white rounded-lg disabled:opacity-50"
          >
            {isSubmitting ? "Creating account…" : "Sign up"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <Link to="/login" className="text-indigo-600 hover:text-indigo-500">
            Log in
          </Link>
        </p>

        <div className="mt-4 text-center space-y-2">
          <p className="text-xs text-gray-500">
            <Link to="/about" className="hover:text-gray-700 transition-colors">
              Learn more about Yappr
            </Link>
          </p>
          <div className="flex justify-center space-x-4 text-xs text-gray-400">
            <Link to="/terms" className="hover:text-gray-600 transition-colors">
              Terms
            </Link>
            <Link
              to="/privacy"
              className="hover:text-gray-600 transition-colors"
            >
              Privacy
            </Link>
          </div>
          <p className="text-xs text-gray-400 mt-3">
            By signing up, you agree to our Terms of Service and Privacy Policy.
          </p>
        </div>
      </div>
    </div>
  );
}

