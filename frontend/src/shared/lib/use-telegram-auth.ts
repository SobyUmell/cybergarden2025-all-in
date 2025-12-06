"use client";

import { retrieveRawInitData } from "@telegram-apps/sdk-react";
import { useEffect, useState } from "react";

const useTelegramAuth = () => {
  const [tma, setTma] = useState<string | undefined>(undefined);
  const initData = retrieveRawInitData();

  useEffect(() => {
    if (initData) {
      setTma(initData);
    }
  });

  return tma;
};

export { useTelegramAuth };
