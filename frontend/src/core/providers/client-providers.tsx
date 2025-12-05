"use client";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/shared/lib";

export const ClientProviders = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </>
  );
};
