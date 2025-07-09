import { useState, useRef, useEffect } from "react";
import { useParams } from "react-router-dom";
import Header from "../components/Header";
import { useAuth } from "../context/AuthContext";
import useChatSocket from "../hooks/useChatSocket";
import MessageBubble from "../components/MessageBubble";

export default function ChatPage() {
  const { roomId = "" } = useParams();
  const { user } = useAuth();
  const { messages, sendMessage } = useChatSocket(roomId);
  const [input, setInput] = useState("");
  const bottomRef = useRef<HTMLDivElement | null>(null);

  /* scroll to newest */
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages.length]);

  function handleSend() {
    const text = input.trim();
    if (!text) return;
    sendMessage(text);
    setInput("");
  }

  function onKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }

  return (
    <div className="h-screen flex flex-col">
      <Header />

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
    </div>
  );
}
