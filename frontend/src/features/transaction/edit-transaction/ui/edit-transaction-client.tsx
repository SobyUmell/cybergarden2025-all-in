"use client";

import { EditTransactionForm } from "@/features/transaction/edit-transaction/ui/edit-transaction-form";
import { useGetTransactionHistory } from "@/features/transaction/hooks";

interface EditTransactionClientProps {
  id: number;
}

export const EditTransactionClient = ({ id }: EditTransactionClientProps) => {
  const { data, isLoading } = useGetTransactionHistory();

  const transaction = data?.data.find((t) => t.id === id);

  if (isLoading) {
    return <p>Загрузка...</p>;
  }

  return <EditTransactionForm transaction={transaction} />;
};
