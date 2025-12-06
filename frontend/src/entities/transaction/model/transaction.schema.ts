import { z } from "zod";

export const transactionSchema = z.object({
  id: z.string(),
  userId: z.string(),
  categoryId: z.string(),
  amount: z.number(),
  type: z.enum(["income", "expense"]),
  description: z.string().optional(),
  date: z.date(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export type Transaction = z.infer<typeof transactionSchema>;

export const createTransactionSchema = transactionSchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export type CreateTransaction = z.infer<typeof createTransactionSchema>;
