import { api } from "./auth";

export interface CheckinResult {
  streak_count: number;
  is_new_checkin: boolean;
}

export interface UserProfile {
  user_id: string;
  daily_streak: number;
  total_checkins: number;
  total_messages: number;
  total_upvotes_received: number;
  can_receive_upvote: boolean;
}

export async function performDailyCheckin(): Promise<CheckinResult> {
  const { data } = await api.post("/api/stats/checkin");
  return data;
}

export async function getUserProfile(userId: string): Promise<UserProfile> {
  const { data } = await api.get(`/api/stats/profile/${userId}`);
  return data;
}

export async function giveUpvote(toUserId: string): Promise<void> {
  await api.post("/api/stats/upvote", { to_user_id: toUserId });
}