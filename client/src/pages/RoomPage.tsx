import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchRooms, createRoom, type Room } from "../api/rooms";
import { performDailyCheckin } from "../api/stats";
import Header from "../components/Header";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import { calculateTimeRemaining, isRoomExpired } from "../utils/timeUtils";

export default function RoomsPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [newName, setNewName] = useState("");
  const [currentTime, setCurrentTime] = useState(new Date());
  const navigate = useNavigate();
  const { user } = useAuth();
  const { showToast } = useToast();

  useEffect(() => {
    refresh();
    // Trigger daily check-in for authenticated users
    if (user && !user.guest) {
      handleDailyCheckin();
    }
  }, [user]);

  // Update current time every minute for time remaining calculations
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 60000); // Update every minute

    return () => clearInterval(timer);
  }, []);

  async function refresh() {
    try {
      const data = await fetchRooms();
      setRooms(data);
    } catch (error: any) {
      console.error("Failed to fetch rooms:", error);
      showToast("Failed to load rooms. Please refresh the page.", "error");
    }
  }

  async function handleDailyCheckin() {
    try {
      const result = await performDailyCheckin();
      if (result.is_new_checkin) {
        const streakMessage = result.streak_count === 1 
          ? "+1 Daily check-in" 
          : `+${result.streak_count} Daily check-in`;
        showToast(streakMessage, "golden", 7000);
        
        // Show achievement notifications
        if (result.new_achievements && result.new_achievements.length > 0) {
          result.new_achievements.forEach((achievement, index) => {
            setTimeout(() => {
              showToast(`🏆 Achievement Unlocked: ${achievement.name}!`, "golden", 8000);
            }, (index + 1) * 2000); // Stagger achievement notifications
          });
        }
      }
    } catch (error: any) {
      // Silently fail check-ins to not interrupt user experience
      console.error("Daily check-in failed:", error);
    }
  }

  async function handleCreate() {
    if (!newName.trim()) {
      showToast("Please enter a room name", "warning");
      return;
    }
    
    try {
      const room = await createRoom(newName.trim());
      setNewName("");
      showToast(`Room "${room.name}" created successfully!`, "success");
      await refresh(); // Refresh to get the latest list from DB
    } catch (error: any) {
      console.error("Failed to create room:", error);
      
      if (error.response?.status === 429) {
        showToast("Maximum number of rooms reached. Please wait for some rooms to expire.", "error");
      } else if (error.response?.status === 409) {
        showToast("You already have an active room. Please wait for it to expire before creating a new one.", "warning");
      } else if (error.response?.status === 400 && error.response?.data?.error?.includes("inappropriate content")) {
        showToast("Room name contains inappropriate content. Please choose a different name.", "warning");
      } else {
        const errorMessage = error.response?.data?.error || "Failed to create room. Please try again.";
        showToast(errorMessage, "error");
      }
    }
  }

  function enterRoom(room: Room) {
    navigate(`/chat/${room.id}`);
  }

  return (
    <div className="h-screen flex flex-col">
      <Header />

      <main className="flex-1 p-6 bg-gray-100">
        {/* Only show room creation for authenticated users */}
        {user && !user.guest && (
          <div className="mb-6 flex gap-2">
            <input
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder="New room name"
              className="flex-1 rounded-md border-gray-300 px-3 py-2 shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            />
            <button
              onClick={handleCreate}
              className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition"
            >
              Create
            </button>
          </div>
        )}

        {/* Message for guest users */}
        {user?.guest && (
          <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              You're logged in as a guest. Sign up to create your own rooms!
            </p>
          </div>
        )}

        {/* room list */}
        <div className="space-y-4">
          {/* Pinned rooms section */}
          {rooms.filter(r => r.is_pinned).length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-700 mb-2">📌 Today's Topics</h2>
              <ul className="space-y-2">
                {rooms.filter(r => r.is_pinned).map((r) => (
                  <li
                    key={r.id}
                    onClick={() => enterRoom(r)}
                    className="cursor-pointer rounded-md bg-gradient-to-r from-indigo-50 to-purple-50 border border-indigo-200 px-4 py-3 shadow hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold text-indigo-900">{r.name}</span>
                          {r.topic_source && (
                            <span className="text-xs bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded-full">
                              {r.topic_source}
                            </span>
                          )}
                        </div>
                        {r.topic_title && (
                          <p className="text-sm text-gray-700 mt-1 line-clamp-2">{r.topic_title}</p>
                        )}
                        {r.topic_description && (
                          <p className="text-xs text-gray-500 mt-1">{r.topic_description}</p>
                        )}
                      </div>
                      <span className="text-indigo-600 text-sm">→</span>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* Regular rooms section */}
          {rooms.filter(r => !r.is_pinned).length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-700 mb-2">Community Rooms</h2>
              <ul className="space-y-2">
                {rooms.filter(r => !r.is_pinned).map((r) => {
                  const timeRemaining = calculateTimeRemaining(r.expires_at);
                  const expired = isRoomExpired(r.expires_at);
                  
                  return (
                    <li
                      key={r.id}
                      onClick={() => !expired && enterRoom(r)}
                      className={`rounded-md bg-white px-4 py-3 shadow ${
                        expired 
                          ? "opacity-50 cursor-not-allowed" 
                          : "cursor-pointer hover:bg-gray-50"
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className={`font-medium ${expired ? "text-gray-400" : "text-gray-800"}`}>
                            {r.name}
                          </span>
                          <span className="text-xs text-gray-500">
                            #{r.id.slice(0, 6)}
                          </span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className={`text-xs font-medium ${
                            expired 
                              ? "text-red-500" 
                              : timeRemaining.includes("Expired")
                              ? "text-red-500"
                              : timeRemaining.includes("m") && !timeRemaining.includes(":")
                              ? "text-orange-500"  // Less than 1 hour
                              : "text-gray-600"
                          }`}>
                            {expired ? "Expired" : timeRemaining}
                          </span>
                          {!expired && (
                            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                          )}
                        </div>
                      </div>
                    </li>
                  );
                })}
              </ul>
            </div>
          )}
          
          {rooms.length === 0 && (
            <p className="text-sm text-gray-500">No rooms yet.</p>
          )}
        </div>
      </main>
    </div>
  );
}
