import { RouterProvider } from "react-router-dom";

import { router } from "@/pages/router";

export function App() {
  return <RouterProvider router={router} />;
}
