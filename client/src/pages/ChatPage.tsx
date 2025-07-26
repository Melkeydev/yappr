import { useState, useRef, useEffect } from "react";
import { useParams } from "react-router-dom";
import Header from "../components/Header";
import { useAuth } from "../context/AuthContext";
import useChatSocket from "../hooks/useChatSocket";
import MessageBubble from "../components/MessageBubble";
import UserProfileModal from "../components/UserProfile";
import { fetchRooms } from "../api/rooms";
import { useToast } from "../context/ToastContext";

export default function ChatPage() {
  const { roomId = "" } = useParams();
  const { user } = useAuth();
  const { messages, sendMessage } = useChatSocket(roomId);
  const [input, setInput] = useState("");
  const [roomInfo, setRoomInfo] = useState<any>(null);
  const [profileModal, setProfileModal] = useState<{
    userId: string;
    username: string;
  } | null>(null);
  const bottomRef = useRef<HTMLDivElement | null>(null);
  const { showToast } = useToast();

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages.length]);

  useEffect(() => {
    async function loadRoomInfo() {
      try {
        const rooms = await fetchRooms();
        const room = rooms.find((r) => r.id === roomId);
        setRoomInfo(room);
      } catch (error) {
        showToast("Failed to load room information", "warning");
      }
    }
    loadRoomInfo();
  }, [roomId, showToast]);

  function handleSend() {
    const text = input.trim();
    if (!text) {
      showToast("Please enter a message", "warning");
      return;
    }

    try {
      sendMessage(text);
      setInput("");
    } catch (error) {
      showToast("Failed to send message. Please try again.", "error");
    }
  }

  function onKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }

  function handleUsernameClick(userId: string, username: string) {
    setProfileModal({ userId, username });
  }

  return (
    <div className="h-screen flex flex-col">
      <Header />

      {/* Topic banner for pinned rooms */}
      {roomInfo?.is_pinned && roomInfo?.topic_title && (
        <div className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-4 py-3 shadow-md">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <h2 className="font-semibold text-lg">
                  {roomInfo.name} - Today's Topic
                </h2>
                <p className="text-sm mt-1 opacity-90">
                  {roomInfo.topic_title}
                </p>
                {roomInfo.topic_description && (
                  <p className="text-xs mt-1 opacity-75">
                    {roomInfo.topic_description}
                  </p>
                )}
              </div>
              {roomInfo.topic_url && (
                <a
                  href={roomInfo.topic_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="ml-4 px-3 py-1 bg-white/20 hover:bg-white/30 rounded-md text-sm transition-colors"
                >
                  View Source →
                </a>
              )}
            </div>
          </div>
        </div>
      )}

      {/* message list */}
      <div className="flex-1 overflow-y-auto bg-gray-50 px-4 py-6 space-y-3">
        {messages.map((m, i) => (
          <div
            key={i}
            className={
              m.username === user?.username
                ? "flex justify-end"
                : "flex justify-start"
            }
          >
            <MessageBubble
              text={m.content}
              mine={m.username === user?.username}
              username={m.username}
              userId={m.user_id}
              onUsernameClick={handleUsernameClick}
            />
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      {/* composer */}
      <div className="p-4 bg-white shadow-inner flex gap-2">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={onKeyDown}
          rows={1}
          placeholder="Type a message…"
          className="flex-1 resize-none rounded-md border-gray-300 px-3 py-2 shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
        />
        <button
          onClick={handleSend}
          className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition disabled:opacity-50"
          disabled={!input.trim()}
        >
          Send
        </button>
      </div>

      {/* User Profile Modal */}
      {profileModal && (
        <UserProfileModal
          userId={profileModal.userId}
          username={profileModal.username}
          isOpen={true}
          onClose={() => setProfileModal(null)}
        />
      )}
    </div>
  );
}
