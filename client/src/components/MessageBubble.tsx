import clsx from "clsx";

type Props = {
  text: string;
  mine: boolean;
  username: string;
  userId?: string;
  timestamp?: string;
  onUsernameClick?: (userId: string, username: string) => void;
};

export default function MessageBubble({ text, mine, username, userId, timestamp, onUsernameClick }: Props) {
  const formatTime = (timestamp?: string) => {
    if (!timestamp) return "";
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div
      className={clsx(
        "max-w-sm px-4 py-2 rounded-lg shadow",
        mine
          ? "bg-indigo-600 text-white self-end rounded-br-none"
          : "bg-white text-gray-800 self-start rounded-bl-none",
      )}
    >
      {!mine && (
        <p className="mb-1 text-xs font-semibold text-indigo-600">
          {userId && onUsernameClick ? (
            <button
              onClick={() => onUsernameClick(userId, username)}
              className="hover:text-indigo-800 hover:underline transition-colors cursor-pointer"
            >
              {username}
            </button>
          ) : (
            username
          )}
        </p>
      )}
      <p className="whitespace-pre-wrap break-words">{text}</p>
      {timestamp && (
        <p className={clsx(
          "text-xs mt-1",
          mine ? "text-indigo-200" : "text-gray-500"
        )}>
          {formatTime(timestamp)}
        </p>
      )}
    </div>
  );
}
