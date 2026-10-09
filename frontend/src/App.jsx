import { useEffect, useState } from "react";

import { authApi } from "./api/client";

import LoginPage from "./pages/LoginPage";
import AdminPage from "./pages/AdminPage";
import WorkshopPage from "./pages/WorkshopPage";
import Layout from "./components/Layout";

import "./App.css";

export default function App() {
  const [user, setUser] = useState(null);
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    async function restoreSession() {
      const token = sessionStorage.getItem("access_token");

      if (!token) {
        setChecking(false);
        return;
      }

      try {
        const response = await authApi.me();
        setUser(response.data);
      } catch {
        sessionStorage.removeItem("access_token");
        setUser(null);
      } finally {
        setChecking(false);
      }
    }

    restoreSession();
  }, []);

  function logout() {
    sessionStorage.removeItem("access_token");
    setUser(null);
  }

  if (checking) {
    return (
      <div className="loading-screen">
        Loading Workshop Desk...
      </div>
    );
  }

  if (!user) {
    return <LoginPage onLogin={setUser} />;
  }

  return (
    <Layout user={user} onLogout={logout}>
      {user.role === "admin" ? (
        <AdminPage />
      ) : user.role === "manager" ||
        user.role === "staff" ? (
        <WorkshopPage user={user} />
      ) : (
        <p>Account role is not supported.</p>
      )}
    </Layout>
  );
}