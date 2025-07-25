import { useEffect, useRef, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";

export type ChatMessage = {
  content: string;
  room_id: string;
  username: string;
  user_id?: string;
};

const WS_URL = import.meta.env.VITE_WEBSOCKET_URL || "ws://localhost:8080";
console.log("this is WS_URL: ", WS_URL);

export default function useChatSocket(roomId: string) {
  const { user } = useAuth();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (!user) return;

    let retries = 0;
    let ws: WebSocket;
    let shouldReconnect = true;

    function connect() {
      ws = new WebSocket(
        `${WS_URL}/ws/joinRoom/${roomId}?userId=${user.id}&username=${encodeURIComponent(
          user.username,
        )}`,
      );

      ws.onopen = () => (retries = 0); // reset back-off on success

      ws.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        setMessages((prev) => [...prev, msg]);
      };

      ws.onclose = (event) => {
        // Check if it's a normal close with error status
        if (event.code === 1008 || event.code === 1003) {
          // Room doesn't exist or expired
          alert(
            "This room has expired or doesn't exist. Redirecting to room list...",
          );
          navigate("/rooms");
          return;
        }

        if (shouldReconnect && retries < 5) {
          retries += 1;
          setTimeout(connect, 500 * retries); // simple back-off
        }
      };

      ws.onerror = () => ws.close(); // close triggers retry
      wsRef.current = ws;
    }

    connect();
    return () => {
      shouldReconnect = false;
      ws.close();
    };
  }, [roomId, user, navigate]);

  function sendMessage(text: string) {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
    }
  }

  return { messages, sendMessage };
}
