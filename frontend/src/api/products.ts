import { api } from "@/api/client";

export interface Product {
  id: string;
  name: string;
  price: number;
  currency_code: string;
}

interface ProductsResponse {
  products: Product[];
}

export async function getProducts(signal?: AbortSignal): Promise<Product[]> {
  const response = await api
    .get("products", { signal })
    .json<ProductsResponse>();

  return response.products;
}
