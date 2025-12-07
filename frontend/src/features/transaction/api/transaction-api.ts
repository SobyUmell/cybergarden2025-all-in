import { queryTma } from "@/shared/api/api-client";
import type {
  Transaction,
  CreateTransaction,
  UpdateTransaction,
  DeleteTransaction,
} from "@/entities/transaction/model/transaction.schema";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

export const transactionApi = {
  getHistory: async (): Promise<{ data: Transaction[] }> => {
    return queryTma<{ data: Transaction[] }>(`${API_BASE}/webapp/datahistory`);
  },

  create: async (transaction: CreateTransaction): Promise<{ status: string }> => {
    return queryTma<{ status: string }>(`${API_BASE}/webapp/addt`, {
      method: "POST",
      body: JSON.stringify(transaction),
    });
  },

  update: async (transaction: UpdateTransaction): Promise<{ status: string }> => {
    return queryTma<{ status: string }>(`${API_BASE}/webapp/updatet`, {
      method: "POST",
      body: JSON.stringify(transaction),
    });
  },

  delete: async (data: DeleteTransaction): Promise<{ status: string }> => {
    return queryTma<{ status: string }>(`${API_BASE}/webapp/deletet`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  chat: async (message: string): Promise<{ response: string }> => {
    return queryTma<{ response: string }>(`${API_BASE}/webapp/chat`, {
      method: "POST",
      body: JSON.stringify({ message }),
    });
  },
};
