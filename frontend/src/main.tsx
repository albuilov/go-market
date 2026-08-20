import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { ConfigProvider, theme } from "antd";

import { App } from "@/app/app";
import { queryClient } from "@/app/query-client";
import "antd/dist/reset.css";
import "@/styles/globals.css";

const root = document.getElementById("root");

if (!root) {
  throw new Error("Root element was not found");
}

createRoot(root).render(
  <StrictMode>
    <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
        token: {
          colorPrimary: "#9e77ed",
          colorBgBase: "#0c111d",
          colorBgContainer: "#161b26",
          colorBorder: "#333741",
          colorTextBase: "#f5f5f6",
          borderRadius: 8,
          fontFamily:
            'Inter, ui-sans-serif, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </ConfigProvider>
  </StrictMode>,
);
