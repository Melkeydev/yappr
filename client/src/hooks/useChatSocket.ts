import { useEffect, useRef, useState } from "react";
import { useAuth } from "../context/AuthContext";

export type ChatMessage = {
  content: string;
  room_id: string;
  username: string;
};

const WS_URL = import.meta.env.VITE_WEBSOCKET_URL || "ws://localhost:8080";

export default function useChatSocket(roomId: string) {
  const { user } = useAuth();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!user) return;

    let retries = 0;
    let ws: WebSocket;

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

      ws.onclose = () => {
        if (retries < 5) {
          retries += 1;
          setTimeout(connect, 500 * retries); // simple back-off
        }
      };

      ws.onerror = () => ws.close(); // close triggers retry
      wsRef.current = ws;
    }

    connect();
    return () => ws.close();
  }, [roomId, user]);

  function sendMessage(text: string) {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(text);
    }
  }

  return { messages, sendMessage };
}
