"use client";

import { DataTable } from "@/shared/ui/data-table";
import { columns } from "../model/transaction.columns";
import { useGetTransactionHistory } from "@/features/transaction/hooks";

export const TransactionTable = () => {
  const { data, isLoading, error } = useGetTransactionHistory();

  if (isLoading) {
    return <div className="p-4">Загрузка транзакций...</div>;
  }

  if (error) {
    return (
      <div className="p-4 text-red-600">
        Ошибка загрузки транзакций: {error.message}
      </div>
    );
  }

  return (
    <DataTable
      columns={columns}
      data={data?.data || []}
      searchKey="description"
      searchPlaceholder="Поиск по описанию..."
    />
  );
};
