"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { z } from "zod"
import { queryTma } from "@/shared/api/api-client"

const createTransactionFormSchema = z.object({
  categoryId: z.string().min(1, "Категория обязательна"),
  amount: z.number().min(0.01, "Сумма должна быть больше 0"),
  type: z.enum(["income", "expense"]),
  description: z.string().optional(),
  date: z.string().min(1, "Дата обязательна"),
})

export type CreateTransactionFormValues = z.infer<typeof createTransactionFormSchema>

const createTransaction = async (data: CreateTransactionFormValues) => {
  return queryTma("/api/transactions", {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export const useCreateTransactionForm = () => {
  const queryClient = useQueryClient()

  const form = useForm<CreateTransactionFormValues>({
    resolver: zodResolver(createTransactionFormSchema),
    defaultValues: {
      categoryId: "",
      amount: 0,
      type: "expense",
      description: "",
      date: new Date().toISOString().split("T")[0],
    },
  })

  const mutation = useMutation({
    mutationFn: createTransaction,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] })
      form.reset()
    },
  })

  const onSubmit = (data: CreateTransactionFormValues) => {
    mutation.mutate(data)
  }

  return {
    form,
    onSubmit: form.handleSubmit(onSubmit),
    isLoading: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    error: mutation.error,
  }
}
