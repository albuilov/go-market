import { useQuery } from "@tanstack/react-query";

import { getProducts } from "@/api/products";

export const productKeys = {
  all: ["products"] as const,
  list: () => [...productKeys.all, "list"] as const,
};

export function useProducts() {
  return useQuery({
    queryKey: productKeys.list(),
    queryFn: ({ signal }) => getProducts(signal),
  });
}
