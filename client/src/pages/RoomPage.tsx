import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchRooms, createRoom, type Room } from "../api/rooms";
import Header from "../components/Header";

export default function RoomsPage() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [newName, setNewName] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    refresh();
  }, []);

  async function refresh() {
    const data = await fetchRooms();
    setRooms(data);
  }

  async function handleCreate() {
    console.log("handleCreate called, newName:", newName);
    if (!newName.trim()) return;
    try {
      console.log("Calling createRoom with:", newName.trim());
      const room = await createRoom(newName.trim());
      console.log("Room created:", room);
      setNewName("");
      await refresh(); // Refresh to get the latest list from DB
    } catch (error: any) {
      console.error("Error creating room:", error);
      if (error.response?.status === 429) {
        alert(
          "Maximum number of rooms reached. Please wait for some rooms to expire.",
        );
      } else {
        alert("Failed to create room. Please try again.");
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
        <div className="mb-6 flex gap-2">
          <input
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="New room name"
            className="flex-1 rounded-md border-gray-300 px-3 py-2 shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
          />
          <button
            onClick={() => {
              console.log("Button clicked!");
              handleCreate();
            }}
            className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 transition"
          >
            Create
          </button>
        </div>

        {/* room list */}
        <ul className="space-y-2">
          {rooms.map((r) => (
            <li
              key={r.id}
              onClick={() => enterRoom(r)}
              className="cursor-pointer rounded-md bg-white px-4 py-3 shadow hover:bg-gray-50"
            >
              <span className="font-medium text-gray-800">{r.name}</span>
              <span className="ml-2 text-xs text-gray-500">
                #{r.id.slice(0, 6)}
              </span>
            </li>
          ))}
          {rooms.length === 0 && (
            <p className="text-sm text-gray-500">No rooms yet.</p>
          )}
        </ul>
      </main>
    </div>
  );
}
