import { Outlet, ScrollRestoration } from "react-router-dom";

import { Footer } from "@/components/footer";
import { Header } from "@/components/header";

export function AppLayout() {
  return (
    <div className="flex min-h-screen flex-col bg-gray-950">
      <Header />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
      <ScrollRestoration />
    </div>
  );
}
