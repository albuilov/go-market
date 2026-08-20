import { LockOutlined } from "@ant-design/icons";
import { Alert, Skeleton } from "antd";
import { Outlet } from "react-router-dom";

import { useUserInfo } from "@/features/auth/use-user-info";

export function ProtectedRoute() {
  const userInfo = useUserInfo();

  if (userInfo.isPending) {
    return (
      <div className="mx-auto max-w-container px-4 py-12 sm:px-6 lg:px-8">
        <Skeleton active paragraph={{ rows: 4 }} />
      </div>
    );
  }

  if (!userInfo.data) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-16 sm:px-6 lg:py-24">
        <Alert
          type="info"
          showIcon
          icon={<LockOutlined />}
          message="Необходимо войти в аккаунт"
          description="После авторизации здесь появятся баланс, история покупок и персональная информация."
        />
      </div>
    );
  }

  return <Outlet />;
}
