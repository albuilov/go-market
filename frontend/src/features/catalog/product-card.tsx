import { ArrowUpOutlined, InboxOutlined } from "@ant-design/icons";
import { Link } from "react-router-dom";

import type { Product } from "@/api/products";
import { formatPrice } from "@/features/catalog/format-price";

const backgrounds = [
  "from-brand-950 to-gray-900 text-brand-400",
  "from-blue-950 to-gray-900 text-blue-400",
  "from-orange-950 to-gray-900 text-orange-400",
] as const;

interface ProductCardProps {
  product: Product;
  index: number;
}

export function ProductCard({ product, index }: ProductCardProps) {
  const background = backgrounds[index % backgrounds.length];

  return (
    <Link
      className="block rounded-2xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-600"
      to={`/product/${encodeURIComponent(product.id)}`}
      aria-label={`${product.name}, открыть товар`}
    >
      <article className="group overflow-hidden rounded-2xl border border-gray-800 bg-gray-900 shadow-xs transition duration-200 hover:-translate-y-0.5 hover:border-gray-700 hover:shadow-lg">
        <div
          className={`flex aspect-[4/3] items-center justify-center bg-gradient-to-br ${background}`}
        >
          <div className="flex size-16 items-center justify-center rounded-xl bg-gray-950/75 shadow-sm ring-1 ring-white/10 backdrop-blur-sm">
            <InboxOutlined className="text-3xl" aria-hidden="true" />
          </div>
        </div>

        <div className="p-4 sm:p-5">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h2 className="text-base font-semibold text-gray-100">{product.name}</h2>
              <p className="mt-2 text-lg font-semibold text-gray-100">
                {formatPrice(product.price, product.currency_code)}
              </p>
            </div>
            <span className="flex size-8 shrink-0 items-center justify-center rounded-full border border-gray-700 text-gray-400 transition group-hover:border-brand-400 group-hover:text-brand-300">
              <ArrowUpOutlined className="rotate-45" aria-hidden="true" />
            </span>
          </div>
        </div>
      </article>
    </Link>
  );
}
