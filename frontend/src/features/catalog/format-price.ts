export function formatPrice(price: number, currencyCode: string): string {
  if (!Number.isFinite(price)) {
    return "Цена недоступна";
  }

  try {
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: currencyCode,
      maximumFractionDigits: 2,
    }).format(price);
  } catch {
    return `${price.toFixed(2)} ${currencyCode}`;
  }
}
