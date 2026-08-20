export function formatPrice(price: number | string, currencyCode: string): string {
  const value = Number(price);

  if (!Number.isFinite(value)) {
    return "Цена недоступна";
  }

  try {
    return new Intl.NumberFormat("ru-RU", {
      style: "currency",
      currency: currencyCode,
      maximumFractionDigits: 2,
    }).format(value);
  } catch {
    return `${value.toFixed(2)} ${currencyCode}`;
  }
}
