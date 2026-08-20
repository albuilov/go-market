export function RouteFallback() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-950 px-4">
      <div className="flex items-center gap-3 text-sm font-medium text-gray-400">
        <span className="size-3 animate-pulse rounded-full bg-brand-600" />
        Загружаем страницу…
      </div>
    </div>
  );
}
