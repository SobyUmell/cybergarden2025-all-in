"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { transactionApi } from "../api";
import type { UpdateTransaction } from "@/entities/transaction/model/transaction.schema";

export const useUpdateTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (transaction: UpdateTransaction) =>
      transactionApi.update(transaction),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
    },
  });
};
