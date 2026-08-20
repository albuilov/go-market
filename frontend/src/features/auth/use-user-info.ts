import { useQuery } from "@tanstack/react-query";

import { getUserInfo } from "@/api/user-info";

export const userKeys = {
  all: ["user"] as const,
  info: () => [...userKeys.all, "info"] as const,
};

export function useUserInfo() {
  return useQuery({
    queryKey: userKeys.info(),
    queryFn: ({ signal }) => getUserInfo(signal),
    retry: false,
    staleTime: 5 * 60_000,
  });
}
