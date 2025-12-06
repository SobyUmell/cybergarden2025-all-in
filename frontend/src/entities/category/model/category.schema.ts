import { z } from "zod";

export const categorySchema = z.object({
  id: z.string(),
  name: z.string().min(1, "Название категории обязательно"),
  description: z.string().optional(),
  icon: z.string().optional(),
  color: z.string().optional(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export type Category = z.infer<typeof categorySchema>;

export const createCategorySchema = categorySchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export type CreateCategory = z.infer<typeof createCategorySchema>;
