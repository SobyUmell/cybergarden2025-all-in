"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { useCreateTransaction } from "../../hooks";

const createTransactionFormSchema = z.object({
  id: z.number().optional(),
  date: z.number(),
  kategoria: z.string().min(1, "Категория обязательна"),
  type: z.string().min(1, "Тип обязателен"),
  amount: z.number().min(0.01, "Сумма должна быть больше 0"),
  description: z.string(),
});

export type CreateTransactionFormValues = z.infer<
  typeof createTransactionFormSchema
>;

export const useCreateTransactionForm = () => {
  const router = useRouter();
  const createMutation = useCreateTransaction();

  const form = useForm<CreateTransactionFormValues>({
    resolver: zodResolver(createTransactionFormSchema),
    defaultValues: {
      id: 0,
      date: Date.now(),
      kategoria: "",
      type: "Списание/Покупка",
      amount: 0,
      description: "",
    },
  });

  const onSubmit = async (data: CreateTransactionFormValues) => {
    try {
      await createMutation.mutateAsync({
        ...data,
        id: data.id ?? 0,
      });
      router.push("/transactions");
    } catch (error) {
      console.error("Failed to create transaction:", error);
    }
  };

  return {
    form,
    onSubmit: form.handleSubmit(onSubmit),
    isLoading: createMutation.isPending,
    isSuccess: createMutation.isSuccess,
    isError: createMutation.isError,
    error: createMutation.error,
  };
};
