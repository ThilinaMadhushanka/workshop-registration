import axios from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 15000,
});

api.interceptors.request.use((config) => {
  const token = sessionStorage.getItem("access_token");

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

export const authApi = {
  login: (data) => api.post("/auth/login", data),
  me: () => api.get("/auth/me"),
  createUser: (data) => api.post("/users", data),
};

export const workshopApi = {
  list: (filters = {}) =>
    api.get("/workshops", { params: filters }),

  create: (data) => api.post("/workshops", data),

  update: (id, data) =>
    api.put(`/workshops/${id}`, data),

  registrations: (id) =>
    api.get(`/workshops/${id}/registrations`),

  register: (id, data) =>
    api.post(`/workshops/${id}/registrations`, data),

  cancel: (id) =>
    api.post(`/registrations/${id}/cancel`),
};

export function getErrorMessage(error) {
  const message = error.response?.data?.error;

  if (typeof message === "string") {
    return message;
  }

  if (!error.response) {
    return "Cannot connect to the server. Please check the backend.";
  }

  return "Something went wrong. Please try again.";
}