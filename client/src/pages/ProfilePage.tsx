import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import Header from "../components/Header";
import { api } from "../api/auth";

type UsernameForm = { username: string };

export default function ProfilePage() {
  const { user, setUser } = useAuth();
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<UsernameForm>({
    defaultValues: {
      username: user?.username || "",
    },
  });

  if (!user || user.guest) {
    navigate("/rooms");
    return null;
  }

  async function onSubmit(values: UsernameForm) {
    try {
      const { data } = await api.put("/api/users/username", { username: values.username });
      
      // Update user in context and localStorage
      const updatedUser = { ...user, username: data.username };
      setUser(updatedUser);
      localStorage.setItem("chat_user", JSON.stringify(updatedUser));
      
      setIsEditing(false);
      alert("Username updated successfully!");
    } catch (error: any) {
      alert(error.response?.data?.error || "Failed to update username");
    }
  }

  return (
    <div className="h-screen flex flex-col">
      <Header />

      <main className="flex-1 p-6 bg-gray-100">
        <div className="max-w-2xl mx-auto">
          <h1 className="text-2xl font-bold mb-6">Profile Settings</h1>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Email
              </label>
              <p className="text-gray-900">{user.email || "Not set"}</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Username
              </label>
              
              {!isEditing ? (
                <div className="flex items-center gap-4">
                  <p className="text-gray-900">{user.username}</p>
                  <button
                    onClick={() => setIsEditing(true)}
                    className="text-sm text-indigo-600 hover:text-indigo-500"
                  >
                    Edit
                  </button>
                </div>
              ) : (
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                  <div>
                    <input
                      type="text"
                      {...register("username", { 
                        required: "Username is required",
                        minLength: {
                          value: 3,
                          message: "Username must be at least 3 characters"
                        },
                        maxLength: {
                          value: 20,
                          message: "Username must be less than 20 characters"
                        },
                        pattern: {
                          value: /^[a-zA-Z0-9_]+$/,
                          message: "Username can only contain letters, numbers, and underscores"
                        }
                      })}
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                    />
                    {errors.username && (
                      <p className="mt-1 text-xs text-red-600">{errors.username.message}</p>
                    )}
                  </div>
                  
                  <div className="flex gap-2">
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className="px-4 py-2 text-sm font-medium bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50"
                    >
                      {isSubmitting ? "Saving..." : "Save"}
                    </button>
                    <button
                      type="button"
                      onClick={() => {
                        setIsEditing(false);
                        reset();
                      }}
                      className="px-4 py-2 text-sm font-medium bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              )}
            </div>
          </div>

          <div className="mt-4">
            <button
              onClick={() => navigate("/rooms")}
              className="text-sm text-gray-600 hover:text-gray-900"
            >
              ← Back to rooms
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}