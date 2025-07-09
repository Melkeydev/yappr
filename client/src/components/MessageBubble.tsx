import clsx from "clsx";

type Props = {
  text: string;
  mine: boolean;
  username: string;
};

export default function MessageBubble({ text, mine, username }: Props) {
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
        <p className="mb-1 text-xs font-semibold text-indigo-600">{username}</p>
      )}
      <p className="whitespace-pre-wrap break-words">{text}</p>
    </div>
  );
}
