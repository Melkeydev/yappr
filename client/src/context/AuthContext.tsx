import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { api } from "../api/auth";

type User = { id: string; username: string; guest?: boolean } | null;
type AuthCtx = {
  user: User;
  setUser: (u: User) => void;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthCtx | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User>(() => {
    const raw = localStorage.getItem("chat_user");
    return raw ? (JSON.parse(raw) as User) : null;
  });

  // keep localStorage in sync if user changes
  useEffect(() => {
    if (user) localStorage.setItem("chat_user", JSON.stringify(user));
    else localStorage.removeItem("chat_user");
  }, [user]);

  async function logout() {
    try {
      await api.get("/api/users/logout"); // clears cookie
    } finally {
      setUser(null);
    }
  }

  return (
    <AuthContext.Provider value={{ user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be inside <AuthProvider>");
  return ctx;
}
