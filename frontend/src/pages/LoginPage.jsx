import { useState } from "react";
import { ArrowRight, LockKeyhole } from "lucide-react";
import { authApi, getErrorMessage } from "../api/client";

export default function LoginPage({ onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await authApi.login({
        email,
        password,
      });

      const token = response.data.token;
      sessionStorage.setItem("access_token", token);

      const userResponse = await authApi.me();

      onLogin(userResponse.data);
    } catch (err) {
      sessionStorage.removeItem("access_token");
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="login-screen">
      <div className="login-brand">
        <div className="brand-mark">W</div>

        <div>
          <h1>Workshop Desk</h1>
          <p>Community Training Centre</p>
        </div>
      </div>

      <div className="login-panel">
        <div className="login-card">
          <div className="login-icon">
            <LockKeyhole size={22} />
          </div>

          <h2>Welcome back</h2>
          <p className="muted">
            Sign in to manage workshops and registrations.
          </p>

          {error && (
            <div className="alert error-alert">{error}</div>
          )}

          <form onSubmit={handleSubmit}>
            <label htmlFor="email">Email address</label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@workshop.local"
              required
            />

            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
            />

            <button
              className="btn btn-primary login-button"
              type="submit"
              disabled={loading}
            >
              {loading ? "Signing in..." : "Sign in"}
              <ArrowRight size={17} />
            </button>
          </form>

          <p className="login-footer">
            Internal staff access only
          </p>
        </div>
      </div>
    </div>
  );
}