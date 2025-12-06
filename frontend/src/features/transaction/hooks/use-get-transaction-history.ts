"use client";

import { useQuery } from "@tanstack/react-query";
import { transactionApi } from "../api";

export const useGetTransactionHistory = () => {
  return useQuery({
    queryKey: ["transactions"],
    queryFn: transactionApi.getHistory,
  });
};
