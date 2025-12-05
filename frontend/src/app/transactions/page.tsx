import type { Metadata } from "next";
import { CategoryTable } from "@/widgets/category-table";

export const metadata: Metadata = {
  title: "История транзакций",
  description: "Просмотр всех транзакций",
};

export default function TransactionsPage() {
  return (
    <div>
      <CategoryTable />
    </div>
  );
}
