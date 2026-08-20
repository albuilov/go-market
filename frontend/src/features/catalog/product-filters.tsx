import { Input, Select } from "antd";

export function ProductFilters() {
  return (
    <div className="sticky top-18 z-30 rounded-2xl border border-gray-800 bg-gray-900/95 p-4 shadow-lg backdrop-blur-sm sm:p-5">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div className="grid gap-3 sm:grid-cols-2 lg:w-2xl">
          <Input disabled placeholder="Поиск по товарам" aria-label="Поиск по товарам" />
          <Select
            disabled
            className="w-full"
            placeholder="Сортировка"
            aria-label="Сортировка товаров"
            options={[]}
          />
        </div>
      </div>
    </div>
  );
}
