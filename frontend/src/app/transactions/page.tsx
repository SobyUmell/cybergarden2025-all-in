import type { Metadata } from "next";
import { TransactionTable } from "@/widgets/transaction-table";
import { Button } from "@/shared/shadcn/ui/button";
import Link from "next/link";

export const metadata: Metadata = {
  title: "История транзакций",
  description: "Просмотр всех транзакций",
};

export default function TransactionsPage() {
  return (
    <div className="space-y-5">
      <Button asChild>
        <Link href="/transactions/create/1">
          Добавить транзакцию
        </Link>
      </Button>
      <TransactionTable />
    </div>
  );
}
