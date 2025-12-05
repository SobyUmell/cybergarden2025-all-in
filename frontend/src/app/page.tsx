import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/shared/shadcn/ui/button";

export const metadata: Metadata = {
  title: "Главная",
  description: "Центр Инвеста - управление финансами",
};

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8">
      <main className="flex w-full max-w-4xl flex-col gap-8">
        <div className="flex flex-col gap-4">
          <h1 className="text-4xl font-bold">Центр Инвеста</h1>
          <p className="text-lg text-muted-foreground">
            Управление финансами и транзакциями
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <Link href="/transactions">
            <Button variant="outline" className="h-32 w-full text-lg">
              История транзакций
            </Button>
          </Link>
          <Link href="/assistant">
            <Button variant="outline" className="h-32 w-full text-lg">
              AI Ассистент
            </Button>
          </Link>
        </div>
      </main>
    </div>
  );
}
