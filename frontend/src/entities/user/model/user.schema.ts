import { z } from "zod";

export const userSchema = z.object({
  id: z.string(),
  email: z.string().email("Некорректный email"),
  name: z.string().min(1, "Имя обязательно"),
  avatar: z.string().url().optional(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export type User = z.infer<typeof userSchema>;

export const createUserSchema = userSchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export type CreateUser = z.infer<typeof createUserSchema>;
