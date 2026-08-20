import { ProductCatalog } from "@/features/catalog/product-catalog";
import { ProductFilters } from "@/features/catalog/product-filters";

export function HomePage() {
  return (
    <>
      <section className="border-b border-gray-800 bg-gray-900/60">
        <div className="mx-auto max-w-container px-4 py-12 sm:px-6 sm:py-16 lg:px-8 lg:py-20">
          <div className="max-w-3xl">
            <h1 className="text-display-sm font-semibold tracking-tight text-gray-100 sm:text-display-md">
              Полезные вещи для рабочего пространства
            </h1>
            <p className="mt-4 max-w-2xl text-lg text-gray-400 sm:text-xl">
              Небольшая коллекция техники, которая помогает сосредоточиться на
              главном и работать с комфортом.
            </p>
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-container px-4 py-8 sm:px-6 sm:py-10 lg:px-8">
        <ProductFilters />
        <ProductCatalog />
      </section>
    </>
  );
}
