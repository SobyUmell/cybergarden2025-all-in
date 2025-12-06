"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { transactionApi } from "../api";
import type { DeleteTransaction } from "@/entities/transaction/model/transaction.schema";

export const useDeleteTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: DeleteTransaction) => transactionApi.delete(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
    },
  });
};
