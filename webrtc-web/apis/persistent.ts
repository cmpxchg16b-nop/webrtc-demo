"use client";

import { useState } from "react";

export interface UsePersistentStorageReturn {
  getValue(): string;
  setValue(value: string): void;
}

export function usePersistentStorage(
  key: PSKey | string,
): UsePersistentStorageReturn {
  const [state, setState] = useState<string | undefined>(undefined);
  if (typeof window === "undefined") {
    return {
      getValue() {
        return "";
      },
      setValue(v: string) {},
    };
  }
  return {
    getValue() {
      return window.localStorage?.getItem(key) ?? state ?? "";
    },
    setValue(value: string) {
      window.localStorage?.setItem(key, value);
      setState(value);
    },
  };
}

export enum PSKey {
  CurrentServer = "current_server",
  PinnedServer = "pinned_server",
  PreferredUsername = "preferred_username",
  HasLoggedIn = "has_logged_in",
  LoggedInAs = "logged_in_as",
  LoggingIn = "logging_in",
  LoggingStartedAt = "logging_in_started_at",
  Unreads = "all_unreads",
}
