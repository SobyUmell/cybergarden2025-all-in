import type { Metadata } from "next";
import { TransactionTable } from "@/widgets/transaction-table";
import { Button } from "@/shared/shadcn/ui/button";
import Link from "next/link";
import { SquarePlus } from "lucide-react";

export const metadata: Metadata = {
  title: "История транзакций",
  description: "Просмотр всех транзакций",
};

export default function TransactionsPage() {
  return (
    <div className="space-y-5">
      <h1 className="text-2xl md:text-3xl font-bold mb-6">Транзакции</h1>
      <Button size={"lg"} asChild>
        <Link href="/transactions/create/1">
          <SquarePlus />
          Добавить транзакцию
        </Link>
      </Button>
      <TransactionTable />
    </div>
  );
}
