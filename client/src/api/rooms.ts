import { api } from "./auth";

export type Room = {
  id: string;
  name: string;
  is_pinned?: boolean;
  created_at: string;
  expires_at: string;
  topic_title?: string;
  topic_description?: string;
  topic_url?: string;
  topic_source?: string;
};

export async function fetchRooms(): Promise<Room[]> {
  try {
    const { data } = await api.get("/ws/getRooms");
    return data;
  } catch (error: any) {
    if (error.response?.status === 404) {
      await new Promise((resolve) => setTimeout(resolve, 100));
      const { data } = await api.get("/ws/getRooms");
      return data;
    }
    throw error;
  }
}

export async function createRoom(name: string): Promise<Room> {
  const body = { name };
  const { data } = await api.post("/ws/createRoom", body);
  return data;
}
