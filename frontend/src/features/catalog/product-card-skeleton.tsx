export function ProductCardSkeleton() {
  return (
    <div className="overflow-hidden rounded-2xl border border-gray-800 bg-gray-900 shadow-xs">
      <div className="aspect-[4/3] animate-pulse bg-gray-800" />
      <div className="space-y-3 p-4 sm:p-5">
        <div className="h-5 w-2/3 animate-pulse rounded bg-gray-800" />
        <div className="h-6 w-1/3 animate-pulse rounded bg-gray-800" />
      </div>
    </div>
  );
}
