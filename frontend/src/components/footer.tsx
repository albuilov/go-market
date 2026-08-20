import { Link } from "react-router-dom";

export function Footer() {
  return (
    <footer className="border-t border-gray-800 bg-gray-950">
      <div className="mx-auto flex max-w-container flex-col gap-4 px-4 py-8 text-sm text-gray-400 sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8">
        <p>© {new Date().getFullYear()} Marketplace</p>
        <div className="flex items-center gap-5">
          <Link className="transition hover:text-brand-700" to="/">
            Каталог
          </Link>
          <Link className="transition hover:text-brand-700" to="/personal">
            Личная страница
          </Link>
        </div>
      </div>
    </footer>
  );
}
