import { useEffect, useState } from "react";
import {
  getUserProfile,
  giveUpvote,
  type UserProfile as UserProfileType,
  type Achievement,
} from "../api/stats";
import { useToast } from "../context/ToastContext";
import { useAuth } from "../context/AuthContext";

interface UserProfileModalProps {
  userId: string;
  username: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function UserProfileModal({
  userId,
  username,
  isOpen,
  onClose,
}: UserProfileModalProps) {
  const [profile, setProfile] = useState<UserProfileType | null>(null);
  const [loading, setLoading] = useState(false);
  const [upvoting, setUpvoting] = useState(false);
  const { showToast } = useToast();
  const { user } = useAuth();

  useEffect(() => {
    if (isOpen && userId) {
      loadProfile();
    }
  }, [isOpen, userId]);

  const loadProfile = async () => {
    try {
      setLoading(true);
      const profileData = await getUserProfile(userId);
      setProfile(profileData);
    } catch (error: any) {
      showToast("Failed to load user profile", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleUpvote = async () => {
    if (!profile || upvoting) return;

    try {
      setUpvoting(true);
      await giveUpvote(userId);
      showToast(`Upvoted ${username}!`, "success");

      await loadProfile();
    } catch (error: any) {
      if (error.response?.status === 400) {
        showToast("Cannot upvote yourself", "warning");
      } else if (error.response?.status === 409) {
        showToast(
          "You've already upvoted this user or used your daily upvote",
          "warning",
        );
      } else {
        showToast("Failed to upvote user", "error");
      }
    } finally {
      setUpvoting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4 relative">
        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600"
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        {/* Header */}
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">{username}</h2>
          <div className="w-16 h-1 bg-indigo-600 rounded"></div>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          </div>
        ) : profile ? (
          <div className="space-y-6">
            {/* Stats Grid */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-gray-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-indigo-600">
                  {profile.daily_streak}
                </div>
                <div className="text-sm text-gray-600">Daily Streak</div>
                <div className="text-xs text-gray-500 mt-1">🔥</div>
              </div>

              <div className="bg-gray-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-green-600">
                  {profile.total_checkins}
                </div>
                <div className="text-sm text-gray-600">Check-ins</div>
                <div className="text-xs text-gray-500 mt-1">📅</div>
              </div>

              <div className="bg-gray-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {profile.total_messages}
                </div>
                <div className="text-sm text-gray-600">Messages</div>
                <div className="text-xs text-gray-500 mt-1">💬</div>
              </div>

              <div className="bg-gray-50 rounded-lg p-4 text-center">
                <div className="text-2xl font-bold text-yellow-600">
                  {profile.total_upvotes_received}
                </div>
                <div className="text-sm text-gray-600">Upvotes</div>
                <div className="text-xs text-gray-500 mt-1">⭐</div>
              </div>
            </div>

            {/* Achievements Section */}
            {profile.achievements && profile.achievements.length > 0 && (
              <div className="pt-4 border-t">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">
                  🏆 Achievements
                </h3>
                <div className="grid grid-cols-1 gap-2">
                  {profile.achievements.map((achievement) => (
                    <div
                      key={achievement.id}
                      className="flex items-center gap-3 bg-gradient-to-r from-yellow-50 to-amber-50 border border-amber-200 rounded-lg p-3"
                    >
                      <span className="text-2xl">{achievement.icon}</span>
                      <div className="flex-1">
                        <div className="font-semibold text-amber-900">
                          {achievement.name}
                        </div>
                        <div className="text-sm text-amber-700">
                          {achievement.description}
                        </div>
                        {achievement.earned_at && (
                          <div className="text-xs text-amber-600 mt-1">
                            Earned:{" "}
                            {new Date(
                              achievement.earned_at,
                            ).toLocaleDateString()}
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Upvote Button */}
            {user && !user.guest && profile.can_receive_upvote && (
              <div className="pt-4 border-t">
                <button
                  onClick={handleUpvote}
                  disabled={upvoting}
                  className="w-full bg-yellow-600 hover:bg-yellow-700 disabled:bg-gray-300 text-white font-medium py-3 px-4 rounded-lg transition-colors flex items-center justify-center gap-2"
                >
                  {upvoting ? (
                    <>
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                      Upvoting...
                    </>
                  ) : (
                    <>
                      <span>⭐</span>
                      Give Upvote
                    </>
                  )}
                </button>
              </div>
            )}

            {/* Upvote status messages */}
            {user && !user.guest && !profile.can_receive_upvote && (
              <div className="pt-4 border-t">
                <div className="text-center text-sm text-gray-500 bg-gray-50 rounded-lg py-3">
                  {userId === user.id
                    ? "You cannot upvote yourself"
                    : "Already upvoted or used daily upvote"}
                </div>
              </div>
            )}

            {user?.guest && (
              <div className="pt-4 border-t">
                <div className="text-center text-sm text-gray-500 bg-yellow-50 rounded-lg py-3">
                  Sign up to give upvotes!
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="text-center py-8 text-gray-500">
            Failed to load profile
          </div>
        )}
      </div>
    </div>
  );
}

