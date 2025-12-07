"use client";

import { useCreateTransactionForm } from "../lib/use-create-transaction-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/shadcn/ui/form";
import { Input } from "@/shared/shadcn/ui/input";
import { Button } from "@/shared/shadcn/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/shadcn/ui/select";
import Link from "next/link";

export const CreateTransactionForm = () => {
  const { form, onSubmit, isLoading } = useCreateTransactionForm();

  return (
    <Form {...form}>
      <form onSubmit={onSubmit} className="space-y-4">
        <Button type="button" variant={"destructive"} asChild>
          <Link href="/transactions">Назад</Link>
        </Button>
        <div className="w-full max-w-[500px] space-y-5">
          <FormField
            control={form.control}
            name="type"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Тип транзакции</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Выберите тип" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="Списание/Покупка">Расход</SelectItem>
                    <SelectItem value="Пополнение/Доход">Доход</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="amount"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Сумма (RUB)</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    placeholder="0.00"
                    step="0.01"
                    min="0"
                    value={field.value || ""}
                    onChange={(e) => {
                      const value = e.target.value;
                      field.onChange(value === "" ? 0 : parseFloat(value));
                    }}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Описание</FormLabel>
                <FormControl>
                  <Input
                    placeholder="Например: Покупка в магазине"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="date"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Дата</FormLabel>
                <FormControl>
                  <Input
                    type="datetime-local"
                    value={
                      field.value
                        ? new Date(field.value)
                            .toISOString()
                            .slice(0, 16)
                        : ""
                    }
                    onChange={(e) => {
                      const value = e.target.value;
                      field.onChange(
                        value ? new Date(value).getTime() : Date.now()
                      );
                    }}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <Button type="submit" className="w-fit mt-2" disabled={isLoading}>
          {isLoading ? "Создание..." : "Создать транзакцию"}
        </Button>
      </form>
    </Form>
  );
};
