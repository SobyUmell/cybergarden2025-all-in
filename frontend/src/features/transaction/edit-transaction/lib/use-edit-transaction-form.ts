"use client";

import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { useUpdateTransaction } from "../../hooks";
import type { Transaction } from "@/entities/transaction/model/transaction.schema";

const editTransactionFormSchema = z.object({
  id: z.number(),
  date: z.number(),
  kategoria: z.string().min(1, "Категория обязательна"),
  type: z.string().min(1, "Тип обязателен"),
  amount: z.number().min(0.01, "Сумма должна быть больше 0"),
  description: z.string(),
});

export type EditTransactionFormValues = z.infer<
  typeof editTransactionFormSchema
>;

export const useEditTransactionForm = (transaction?: Transaction) => {
  const router = useRouter();
  const updateMutation = useUpdateTransaction();

  const form = useForm<EditTransactionFormValues>({
    resolver: zodResolver(editTransactionFormSchema),
    defaultValues: {
      id: 0,
      date: Date.now(),
      kategoria: "",
      type: "Списание/Покупка",
      amount: 0,
      description: "",
    },
  });

  useEffect(() => {
    if (transaction) {
      form.reset({
        id: transaction.id,
        date: transaction.date,
        kategoria: transaction.kategoria,
        type: transaction.type,
        amount: transaction.amount,
        description: transaction.description,
      });
    }
  }, [transaction, form]);

  const onSubmit = async (data: EditTransactionFormValues) => {
    try {
      await updateMutation.mutateAsync(data);
      router.push("/transactions");
    } catch (error) {
      console.error("Failed to update transaction:", error);
    }
  };

  return {
    form,
    onSubmit: form.handleSubmit(onSubmit),
    isLoading: updateMutation.isPending,
    isSuccess: updateMutation.isSuccess,
    isError: updateMutation.isError,
    error: updateMutation.error,
  };
};
