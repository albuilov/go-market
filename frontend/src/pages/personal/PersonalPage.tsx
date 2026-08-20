const sections = [
  {
    title: "Баланс кошелька",
    description: "Здесь будет отображаться текущий баланс и доступные средства.",
  },
  {
    title: "История покупок",
    description: "Здесь появится список заказов и приобретённых товаров.",
  },
  {
    title: "Персональная информация",
    description: "Здесь можно будет просматривать и изменять данные профиля.",
  },
] as const;

export function PersonalPage() {
  return (
    <section className="mx-auto max-w-container px-4 py-12 sm:px-6 sm:py-16 lg:px-8">
      <div className="max-w-2xl">
        <p className="text-sm font-semibold text-brand-700">Личный кабинет</p>
        <h1 className="mt-3 text-display-sm font-semibold tracking-tight text-gray-100 sm:text-display-md">
          Персональная страница
        </h1>
        <p className="mt-4 text-lg text-gray-400">
          Управление данными аккаунта и покупками.
        </p>
      </div>

      <div className="mt-10 grid gap-6 md:grid-cols-3">
        {sections.map((section) => (
          <article
            className="rounded-2xl border border-gray-800 bg-gray-900 p-6 shadow-xs"
            key={section.title}
          >
            <h2 className="text-lg font-semibold text-gray-100">{section.title}</h2>
            <p className="mt-3 text-sm leading-6 text-gray-400">
              {section.description}
            </p>
          </article>
        ))}
      </div>
    </section>
  );
}
