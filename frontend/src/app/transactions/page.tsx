import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/shared/shadcn/ui/button";

export const metadata: Metadata = {
  title: "История транзакций",
  description: "Просмотр всех транзакций",
};

export default function TransactionsPage() {
  return (
    <div className="flex min-h-screen flex-col p-8">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex flex-col gap-2">
          <h1 className="text-3xl font-bold">История транзакций</h1>
          <p className="text-muted-foreground">
            Просмотр и управление транзакциями
          </p>
        </div>
        <Link href="/">
          <Button variant="outline">На главную</Button>
        </Link>
      </div>

      <div className="flex flex-col gap-4">
        <p className="text-muted-foreground">
          Здесь будет отображаться список транзакций
        </p>
      </div>
    </div>
  );
}
