import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button } from "antd";
import { useNavigate, useParams } from "react-router-dom";

export function ProductPage() {
  const { uuid } = useParams<{ uuid: string }>();
  const navigate = useNavigate();

  return (
    <section className="mx-auto max-w-container px-4 py-12 sm:px-6 sm:py-16 lg:px-8">
      <Button
        type="text"
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate("/")}
      >
        Вернуться в каталог
      </Button>

      <div className="mt-8 overflow-hidden rounded-2xl border border-gray-800 bg-gray-900 shadow-xs">
        <div className="grid lg:grid-cols-2">
          <div className="flex min-h-72 items-center justify-center bg-gradient-to-br from-brand-950 to-gray-900 p-8 sm:min-h-96">
            <div className="rounded-2xl bg-gray-950/80 px-8 py-6 text-center shadow-sm ring-1 ring-white/10">
              <p className="text-sm font-medium text-brand-300">Изображение товара</p>
            </div>
          </div>

          <div className="p-6 sm:p-10 lg:p-12">
            <p className="text-sm font-semibold text-brand-700">Карточка товара</p>
            <h1 className="mt-3 text-display-sm font-semibold tracking-tight text-gray-100">
              Детальная информация
            </h1>
            <p className="mt-4 text-base leading-7 text-gray-400">
              Позже здесь появятся описание, стоимость и действия с товаром.
            </p>

            <div className="mt-8 rounded-xl border border-gray-800 bg-gray-950 p-4">
              <p className="text-xs font-semibold uppercase tracking-wide text-gray-400">
                UUID товара
              </p>
              <p className="mt-2 break-all font-mono text-sm text-gray-100">
                {uuid ?? "Не указан"}
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
