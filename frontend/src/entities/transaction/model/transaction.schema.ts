import { z } from "zod";

export const transactionSchema = z.object({
  id: z.number(),
  date: z.number(),
  kategoria: z.string(),
  type: z.string(),
  amount: z.number(),
  description: z.string(),
});

export type Transaction = z.infer<typeof transactionSchema>;

export const createTransactionSchema = transactionSchema.omit({
  id: true,
});

export type CreateTransaction = z.infer<typeof createTransactionSchema>;

export const updateTransactionSchema = transactionSchema;

export type UpdateTransaction = z.infer<typeof updateTransactionSchema>;

export const deleteTransactionSchema = z.object({
  id: z.number(),
});

export type DeleteTransaction = z.infer<typeof deleteTransactionSchema>;
