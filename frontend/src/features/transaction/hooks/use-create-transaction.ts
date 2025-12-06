"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { transactionApi } from "../api";
import type { CreateTransaction } from "@/entities/transaction/model/transaction.schema";

export const useCreateTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (transaction: CreateTransaction) =>
      transactionApi.create(transaction),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
    },
  });
};
