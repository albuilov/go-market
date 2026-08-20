import { createBrowserRouter } from "react-router-dom";

import { RouteFallback } from "@/components/route-fallback";
import { AppLayout } from "@/layouts/AppLayout";

export const router = createBrowserRouter([
  {
    Component: AppLayout,
    HydrateFallback: RouteFallback,
    children: [
      {
        index: true,
        lazy: async () => {
          const { HomePage } = await import("@/pages/home/HomePage");
          return { Component: HomePage };
        },
      },
      {
        path: "product/:uuid",
        lazy: async () => {
          const { ProductPage } = await import("@/pages/product/ProductPage");
          return { Component: ProductPage };
        },
      },
      {
        lazy: async () => {
          const { ProtectedRoute } = await import(
            "@/components/auth/ProtectedRoute"
          );
          return { Component: ProtectedRoute };
        },
        children: [
          {
            path: "personal",
            lazy: async () => {
              const { PersonalPage } = await import(
                "@/pages/personal/PersonalPage"
              );
              return { Component: PersonalPage };
            },
          },
        ],
      },
    ],
  },
]);
