import { useEffect, useState } from "react";
import {
  CalendarDays,
  Users,
  CheckCircle2,
  Plus,
  RefreshCw,
  X,
  Edit3,
} from "lucide-react";

import {
  workshopApi,
  getErrorMessage,
} from "../api/client";

const emptyWorkshop = {
  code: "",
  title: "",
  instructor: "",
  startsAt: "",
  capacity: 20,
  status: "scheduled",
};

function toInputDate(value) {
  if (!value) return "";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "";

  const local = new Date(
    date.getTime() - date.getTimezoneOffset() * 60000
  );

  return local.toISOString().slice(0, 16);
}

export default function WorkshopPage({ user }) {
  const isManager = user.role === "manager";

  const [workshops, setWorkshops] = useState([]);
  const [selected, setSelected] = useState(null);
  const [registrations, setRegistrations] = useState([]);

  const [filters, setFilters] = useState({
    status: "",
    from: "",
    to: "",
    availableOnly: false,
  });

  const [showWorkshopForm, setShowWorkshopForm] =
    useState(false);

  const [editingID, setEditingID] = useState(null);

  const [workshopForm, setWorkshopForm] =
    useState(emptyWorkshop);

  const [attendee, setAttendee] = useState({
    attendeeName: "",
    attendeeEmail: "",
  });

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  async function loadWorkshops() {
    setLoading(true);
    setError("");

    try {
      const params = {};

      if (filters.status) params.status = filters.status;
      if (filters.from) params.from = filters.from;
      if (filters.to) params.to = filters.to;

      if (filters.availableOnly) {
        params.availableOnly = true;
      }

      const response = await workshopApi.list(params);
      setWorkshops(response.data);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadWorkshops();
  }, [
    filters.status,
    filters.from,
    filters.to,
    filters.availableOnly,
  ]);

  async function openRegistrations(workshop) {
    setSelected(workshop);
    setRegistrations([]);
    setError("");

    try {
      const response = await workshopApi.registrations(
        workshop.id
      );

      setRegistrations(response.data);
    } catch (err) {
      setError(getErrorMessage(err));
    }
  }

  function openCreateForm() {
    setError("");
    setNotice("");
    setEditingID(null);
    setWorkshopForm({ ...emptyWorkshop });
    setShowWorkshopForm(true);
  }

  function openEditForm(workshop) {
    setEditingID(workshop.id);

    setWorkshopForm({
      code: workshop.code,
      title: workshop.title,
      instructor: workshop.instructor,
      startsAt: toInputDate(workshop.startsAt),
      capacity: workshop.capacity,
      status: workshop.status,
    });

    setShowWorkshopForm(true);
  }

  function updateWorkshopField(e) {
    const { name, value } = e.target;

    setWorkshopForm((prev) => ({
      ...prev,
      [name]: name === "capacity" ? Number(value) : value,
    }));
  }

  async function saveWorkshop(e) {
    e.preventDefault();
    setSaving(true);
    setError("");
    setNotice("");

    try {
      const date = new Date(workshopForm.startsAt);

      if (Number.isNaN(date.getTime())) {
        throw new Error("Invalid workshop date");
      }

      const payload = {
        ...workshopForm,
        startsAt: date.toISOString(),
      };

      if (editingID) {
        await workshopApi.update(editingID, payload);
      } else {
        await workshopApi.create(payload);
      }

      setShowWorkshopForm(false);
      setEditingID(null);
      setNotice("Workshop saved successfully.");

      await loadWorkshops();
    } catch (err) {
      setError(
        err.response ? getErrorMessage(err) : err.message
      );
    } finally {
      setSaving(false);
    }
  }

  async function registerAttendee(e) {
    e.preventDefault();

    if (!selected) return;

    setSaving(true);
    setError("");
    setNotice("");

    try {
      await workshopApi.register(selected.id, attendee);

      setAttendee({
        attendeeName: "",
        attendeeEmail: "",
      });

      setNotice("Attendee registered successfully.");

      await refreshSelected();
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  }

  async function refreshSelected() {
    if (!selected) return;

    const [workshopResult, registrationResult] =
      await Promise.all([
        workshopApi.list(),
        workshopApi.registrations(selected.id),
      ]);

    const latest = workshopResult.data.find(
      (w) => w.id === selected.id
    );

    if (latest) setSelected(latest);

    setRegistrations(registrationResult.data);
    await loadWorkshops();
  }

  async function cancelRegistration(id) {
    if (!window.confirm("Cancel this registration?")) {
      return;
    }

    setSaving(true);
    setError("");
    setNotice("");

    try {
      await workshopApi.cancel(id);
      setNotice("Registration cancelled successfully.");
      await refreshSelected();
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setSaving(false);
    }
  }

  const total = workshops.length;

  const scheduled = workshops.filter(
    (w) => w.status === "scheduled"
  ).length;

  const seatsAvailable = workshops.reduce(
    (sum, w) => sum + w.availableSeats,
    0
  );

  return (
    <>
      <div className="page-heading">
        <div>
          <div className="eyebrow">WORKSHOP OPERATIONS</div>
          <h1>Workshop overview</h1>
          <p>
            Manage scheduled workshops, attendance and
            registrations.
          </p>
        </div>

        {isManager && (
          <button
            className="btn btn-primary"
            onClick={openCreateForm}
          >
            <Plus size={17} />
            New workshop
          </button>
        )}
      </div>

      {error && (
        <div className="alert error-alert">{error}</div>
      )}

      {notice && (
        <div className="alert success-alert">{notice}</div>
      )}

      <div className="stats-grid">
        <StatCard
          icon={<CalendarDays size={20} />}
          label="Total workshops"
          value={total}
        />

        <StatCard
          icon={<CheckCircle2 size={20} />}
          label="Scheduled"
          value={scheduled}
        />

        <StatCard
          icon={<Users size={20} />}
          label="Available seats"
          value={seatsAvailable}
        />
      </div>

      <section className="panel">
        <div className="panel-heading between">
          <div>
            <h2>Workshops</h2>
            <p className="muted">
              Browse upcoming and past sessions
            </p>
          </div>

          <button
            className="btn btn-outline"
            onClick={loadWorkshops}
          >
            <RefreshCw size={16} />
            Refresh
          </button>
        </div>

        <div className="filter-bar">
          <select
            value={filters.status}
            onChange={(e) =>
              setFilters({
                ...filters,
                status: e.target.value,
              })
            }
          >
            <option value="">All statuses</option>
            <option value="scheduled">Scheduled</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>

          <input
            type="date"
            aria-label="From date"
            value={filters.from}
            onChange={(e) =>
              setFilters({
                ...filters,
                from: e.target.value,
              })
            }
          />

          <input
            type="date"
            aria-label="To date"
            value={filters.to}
            onChange={(e) =>
              setFilters({
                ...filters,
                to: e.target.value,
              })
            }
          />

          <label className="filter-checkbox">
            <input
              type="checkbox"
              checked={filters.availableOnly}
              onChange={(e) =>
                setFilters({
                  ...filters,
                  availableOnly: e.target.checked,
                })
              }
            />
            Seats available
          </label>
        </div>

        <div className="table-scroll">
          <table className="data-table">
            <thead>
              <tr>
                <th>Workshop</th>
                <th>Date & time</th>
                <th>Seats</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {loading ? (
                <tr>
                  <td colSpan="5">Loading workshops...</td>
                </tr>
              ) : workshops.length === 0 ? (
                <tr>
                  <td colSpan="5">
                    No workshops match your filters.
                  </td>
                </tr>
              ) : (
                workshops.map((w) => (
                  <tr key={w.id}>
                    <td>
                      <strong>{w.title}</strong>
                      <small>
                        {w.code} · {w.instructor}
                      </small>
                    </td>

                    <td>
                      {new Date(w.startsAt).toLocaleString()}
                    </td>

                    <td>
                      <strong>
                        {w.capacity - w.availableSeats}
                        /{w.capacity}
                      </strong>
                      <small>
                        {w.availableSeats} remaining
                      </small>
                    </td>

                    <td>
                      <span className={`status ${w.status}`}>
                        {w.status}
                      </span>
                    </td>

                    <td>
                      <div className="table-actions">
                        <button
                          className="btn btn-outline btn-small"
                          onClick={() => openRegistrations(w)}
                        >
                          View
                        </button>

                        {isManager && (
                          <button
                            className="icon-button"
                            aria-label="Edit workshop"
                            onClick={() => openEditForm(w)}
                          >
                            <Edit3 size={16} />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>

      {showWorkshopForm && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>
                {editingID ? "Edit workshop" : "New workshop"}
              </h2>

              <button
                className="icon-button"
                aria-label="Close"
                onClick={() => setShowWorkshopForm(false)}
              >
                <X size={20} />
              </button>
            </div>

            {error && (
                <div className="alert error-alert">{error}</div>
            )}

            {notice && (
                <div className="alert success-alert">{notice}</div>
            )}

            <form onSubmit={saveWorkshop}>
              <label>Workshop code</label>
              <input
                name="code"
                value={workshopForm.code}
                onChange={updateWorkshopField}
                maxLength={50}
                required
              />

              <label>Workshop title</label>
              <input
                name="title"
                value={workshopForm.title}
                onChange={updateWorkshopField}
                maxLength={200}
                required
              />

              <label>Instructor</label>
              <input
                name="instructor"
                value={workshopForm.instructor}
                onChange={updateWorkshopField}
                required
              />

              <label>Date & time</label>
              <input
                type="datetime-local"
                name="startsAt"
                value={workshopForm.startsAt}
                onChange={updateWorkshopField}
                required
              />

              <div className="form-columns">
                <div>
                  <label>Capacity</label>
                  <input
                    name="capacity"
                    type="number"
                    min={1}
                    value={workshopForm.capacity}
                    onChange={updateWorkshopField}
                    required
                  />
                </div>

                <div>
                  <label>Status</label>
                  <select
                    name="status"
                    value={workshopForm.status}
                    onChange={updateWorkshopField}
                  >
                    <option value="scheduled">Scheduled</option>
                    <option value="completed">Completed</option>
                    <option value="cancelled">Cancelled</option>
                  </select>
                </div>
              </div>

              <div className="modal-actions">
                <button
                  className="btn btn-outline"
                  type="button"
                  onClick={() => setShowWorkshopForm(false)}
                >
                  Cancel
                </button>

                <button
                  className="btn btn-primary"
                  type="submit"
                  disabled={saving}
                >
                  {saving ? "Saving..." : "Save workshop"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {selected && (
        <div className="modal-overlay">
          <div className="modal modal-wide">
            <div className="modal-header">
              <div>
                <h2>{selected.title}</h2>
                <p className="muted">
                  {selected.code} · {selected.availableSeats}
                  {" "}seats remaining
                </p>
              </div>

              <button
                className="icon-button"
                aria-label="Close"
                onClick={() => setSelected(null)}
              >
                <X size={20} />
              </button>
            </div>

            {error && (
                <div className="alert error-alert">{error}</div>
            )}

            {notice && (
                <div className="alert success-alert">{notice}</div>
            )}

            <h3>Register attendee</h3>

            <form
              className="registration-form"
              onSubmit={registerAttendee}
            >
              <input
                placeholder="Attendee name"
                aria-label="Attendee name"
                value={attendee.attendeeName}
                onChange={(e) =>
                  setAttendee({
                    ...attendee,
                    attendeeName: e.target.value,
                  })
                }
                required
              />

              <input
                type="email"
                placeholder="Email address"
                aria-label="Attendee email"
                value={attendee.attendeeEmail}
                onChange={(e) =>
                  setAttendee({
                    ...attendee,
                    attendeeEmail: e.target.value,
                  })
                }
                required
              />

              <button
                className="btn btn-primary"
                type="submit"
                disabled={
                  saving ||
                  selected.availableSeats <= 0 ||
                  selected.status !== "scheduled"
                }
              >
                Register
              </button>
            </form>

            <h3>Registration history</h3>

            <div className="table-scroll">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Attendee</th>
                    <th>Registered</th>
                    <th>Status</th>
                    <th>Action</th>
                  </tr>
                </thead>

                <tbody>
                  {registrations.length === 0 ? (
                    <tr>
                      <td colSpan="4">
                        No registrations recorded.
                      </td>
                    </tr>
                  ) : (
                    registrations.map((r) => (
                      <tr key={r.id}>
                        <td>
                          <strong>{r.attendeeName}</strong>
                          <small>{r.attendeeEmail}</small>
                        </td>

                        <td>
                          {new Date(
                            r.registeredAt
                          ).toLocaleString()}
                          <small>User ID: {r.registeredBy}</small>
                        </td>

                        <td>
                          <span className={`status ${r.status}`}>
                            {r.status}
                          </span>

                          {r.cancelledAt && (
                            <small>
                              Cancelled{" "}
                              {new Date(
                                r.cancelledAt
                              ).toLocaleString()}
                              {" "}by User {r.cancelledBy}
                            </small>
                          )}
                        </td>

                        <td>
                          {r.status === "active" && (
                            <button
                              className="btn btn-danger btn-small"
                              disabled={saving}
                              onClick={() =>
                                cancelRegistration(r.id)
                              }
                            >
                              Cancel
                            </button>
                          )}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            <div className="modal-actions">
              <button
                className="btn btn-outline"
                onClick={() => setSelected(null)}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

function StatCard({ icon, label, value }) {
  return (
    <div className="stat-card">
      <div className="stat-top">
        <span>{label}</span>
        {icon}
      </div>

      <strong>{value}</strong>
    </div>
  );
}