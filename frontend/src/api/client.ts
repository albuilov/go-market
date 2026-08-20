import ky from "ky";

export const api = ky.create({
  prefix: "/api/v1",
  timeout: 10_000,
  retry: {
    limit: 1,
  },
  credentials: "same-origin",
  headers: {
    Accept: "application/json",
  },
});
