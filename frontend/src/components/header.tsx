import { ShoppingOutlined, UserOutlined } from "@ant-design/icons";
import { Avatar, Button, Skeleton } from "antd";
import { Link, NavLink } from "react-router-dom";

import { useUserInfo } from "@/features/auth/use-user-info";

export function Header() {
  const userInfo = useUserInfo();

  return (
    <header className="sticky top-0 z-50 border-b border-gray-800 bg-gray-950/95 backdrop-blur-sm">
      <div className="mx-auto flex h-18 max-w-container items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
        <Link
          className="flex shrink-0 items-center gap-3"
          to="/"
          aria-label="Marketplace, главная"
        >
          <span className="flex size-10 items-center justify-center rounded-xl bg-brand-600 text-white shadow-xs">
            <ShoppingOutlined className="text-xl" aria-hidden="true" />
          </span>
          <span className="text-lg font-semibold tracking-tight text-gray-100">
            Marketplace
          </span>
        </Link>

        <div className="flex items-center gap-3">
          <NavLink
            className={({ isActive }) =>
              `hidden rounded-full px-3 py-1 text-sm font-medium transition sm:block ${
                isActive
                  ? "bg-brand-950 text-brand-300"
                  : "text-gray-400 hover:bg-gray-800 hover:text-gray-100"
              }`
            }
            to="/"
            end
          >
            Каталог
          </NavLink>

          {userInfo.isPending && <Skeleton.Avatar active size="default" />}

          {userInfo.isSuccess && userInfo.data && (
            <Link to="/personal" aria-label="Личная страница">
              <Avatar
                className="bg-brand-600"
                src={userInfo.data.avatar_url}
                icon={!userInfo.data.name ? <UserOutlined /> : undefined}
              >
                {getInitials(userInfo.data.name)}
              </Avatar>
            </Link>
          )}

          {(userInfo.isError || (userInfo.isSuccess && !userInfo.data)) && (
            <Button type="primary">Войти</Button>
          )}
        </div>
      </div>
    </header>
  );
}

function getInitials(name: string): string {
  return name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}
