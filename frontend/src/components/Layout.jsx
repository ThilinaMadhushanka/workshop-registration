import {
  LayoutDashboard,
  Users,
  LogOut,
  CalendarDays,
} from "lucide-react";

export default function Layout({
  user,
  children,
  onLogout,
}) {
  const isAdmin = user.role === "admin";

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="sidebar-brand">
          <div className="sidebar-logo">W</div>

          <div>
            <strong>Workshop Desk</strong>
            <small>Management Portal</small>
          </div>
        </div>

        <div className="sidebar-section">WORKSPACE</div>

        <nav className="side-nav">
          <div className="nav-item active">
            {isAdmin ? (
              <Users size={18} />
            ) : (
              <LayoutDashboard size={18} />
            )}

            <span>
              {isAdmin ? "User Management" : "Overview"}
            </span>
          </div>

          {!isAdmin && (
            <div className="nav-item">
              <CalendarDays size={18} />
              <span>Workshop Operations</span>
            </div>
          )}
        </nav>

        <div className="sidebar-bottom">
          <div className="current-user">
            <div className="avatar">
              {user.role[0].toUpperCase()}
            </div>

            <div>
              <strong>{user.role}</strong>
              <small>Staff account</small>
            </div>
          </div>

          <button
            className="logout-button"
            onClick={onLogout}
          >
            <LogOut size={17} />
            Sign out
          </button>
        </div>
      </aside>

      <div className="main-area">
        <header className="topbar">
          <span>Community Training Centre</span>

          <div className="topbar-user">
            <span className="online-indicator" />
            {user.role.toUpperCase()}
          </div>
        </header>

        <main className="page-content">{children}</main>
      </div>
    </div>
  );
}