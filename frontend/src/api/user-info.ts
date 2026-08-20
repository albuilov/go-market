import { HTTPError } from "ky";

import { api } from "@/api/client";

export interface UserInfo {
  id: string;
  name: string;
  avatar_url?: string;
}

export async function getUserInfo(
  signal?: AbortSignal,
): Promise<UserInfo | null> {
  try {
    return await api.get("user-info", { signal }).json<UserInfo>();
  } catch (error) {
    if (
      error instanceof HTTPError &&
      (error.response.status === 401 || error.response.status === 404)
    ) {
      return null;
    }

    throw error;
  }
}
