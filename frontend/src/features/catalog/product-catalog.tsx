import { InboxOutlined, ReloadOutlined } from "@ant-design/icons";
import { Button } from "antd";

import { ProductCard } from "@/features/catalog/product-card";
import { ProductCardSkeleton } from "@/features/catalog/product-card-skeleton";
import { useProducts } from "@/features/catalog/use-products";

export function ProductCatalog() {
  const products = useProducts();

  return (
    <div className="mt-10 pb-4 sm:mt-12 sm:pb-6">
      {products.isPending && (
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4" aria-label="Загрузка товаров">
          {Array.from({ length: 4 }, (_, index) => (
            <ProductCardSkeleton key={index} />
          ))}
        </div>
      )}

      {products.isError && (
        <div className="flex min-h-72 flex-col items-center justify-center rounded-2xl border border-red-900 bg-red-950/40 px-6 py-12 text-center">
          <span className="flex size-12 items-center justify-center rounded-full bg-red-900/70 text-red-300">
            <ReloadOutlined className="text-xl" aria-hidden="true" />
          </span>
          <h3 className="mt-4 text-lg font-semibold text-gray-100">
            Не удалось загрузить товары
          </h3>
          <p className="mt-2 max-w-md text-sm text-gray-400">
            Проверьте, что backend запущен, и попробуйте выполнить запрос ещё раз.
          </p>
          <Button
            className="mt-6"
            type="primary"
            icon={<ReloadOutlined />}
            onClick={() => void products.refetch()}
          >
            Попробовать снова
          </Button>
        </div>
      )}

      {products.isSuccess && products.data.length === 0 && (
        <div className="flex min-h-72 flex-col items-center justify-center rounded-2xl border border-gray-800 bg-gray-900 px-6 py-12 text-center">
          <span className="flex size-12 items-center justify-center rounded-full bg-gray-800 text-gray-400">
            <InboxOutlined className="text-xl" aria-hidden="true" />
          </span>
          <h3 className="mt-4 text-lg font-semibold text-gray-100">
            Каталог пока пуст
          </h3>
          <p className="mt-2 text-sm text-gray-400">
            Новые товары появятся здесь немного позже.
          </p>
        </div>
      )}

      {products.isSuccess && products.data.length > 0 && (
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {products.data.map((product, index) => (
            <ProductCard key={product.id} product={product} index={index} />
          ))}
        </div>
      )}
    </div>
  );
}
