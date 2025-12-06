import { retrieveRawInitData } from "@telegram-apps/sdk-react";

export const queryTma = async <T>(
  url: string,
  options: RequestInit = {},
): Promise<T> => {
  const tma = retrieveRawInitData();

  const res = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: tma ? `tma ${tma}` : "",
      ...(options.headers || {}),
    },
  });

  if (!res.ok) {
    throw new Error(`Error: ${res.status}`);
  }

  return res.json() as Promise<T>;
}
