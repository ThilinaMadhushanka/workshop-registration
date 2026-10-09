import { useState } from "react";
import { UserPlus, ShieldCheck } from "lucide-react";
import { authApi, getErrorMessage } from "../api/client";

const initialForm = {
  name: "",
  email: "",
  password: "",
  role: "staff",
};

export default function AdminPage() {
  const [form, setForm] = useState(initialForm);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  function handleChange(e) {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  }

  async function handleSubmit(e) {
    e.preventDefault();

    setLoading(true);
    setError("");
    setMessage("");

    try {
      await authApi.createUser(form);

      setMessage("User account created successfully.");
      setForm(initialForm);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <div className="page-heading">
        <div>
          <div className="eyebrow">ADMINISTRATION</div>
          <h1>User management</h1>
          <p>Create accounts for centre managers and staff.</p>
        </div>
      </div>

      <div className="admin-grid">
        <section className="panel">
          <div className="panel-heading">
            <UserPlus size={20} />
            <h2>Create user account</h2>
          </div>

          {message && (
            <div className="alert success-alert">{message}</div>
          )}

          {error && (
            <div className="alert error-alert">{error}</div>
          )}

          <form onSubmit={handleSubmit}>
            <label>Full name</label>
            <input
              name="name"
              value={form.name}
              onChange={handleChange}
              placeholder="Enter full name"
              required
            />

            <label>Email address</label>
            <input
              name="email"
              type="email"
              value={form.email}
              onChange={handleChange}
              placeholder="name@example.com"
              required
            />

            <label>Temporary password</label>
            <input
              name="password"
              type="password"
              value={form.password}
              onChange={handleChange}
              minLength={8}
              required
            />

            <label>Account role</label>
            <select
              name="role"
              value={form.role}
              onChange={handleChange}
            >
              <option value="staff">Staff</option>
              <option value="manager">Manager</option>
            </select>

            <div className="form-footer">
              <button
                type="submit"
                className="btn btn-primary"
                disabled={loading}
              >
                {loading ? "Creating..." : "Create account"}
              </button>
            </div>
          </form>
        </section>

        <aside className="panel info-panel">
          <ShieldCheck size={23} />

          <h2>Account permissions</h2>

          <p>
            New accounts are assigned permissions
            based on their selected role.
          </p>

          <div className="permission-item">
            <strong>Manager</strong>
            <span>
              Create and edit workshops, manage registrations
              and view history.
            </span>
          </div>

          <div className="permission-item">
            <strong>Staff</strong>
            <span>
              View workshops, register attendees and
              manage cancellations.
            </span>
          </div>

          <p className="muted">
            Only administrators can create user accounts.
          </p>
        </aside>
      </div>
    </>
  );
}