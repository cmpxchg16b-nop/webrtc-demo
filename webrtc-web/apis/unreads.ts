import { useState } from "react";

const UNREADS_STORAGE_KEY = "webrtc_unread_message_ids";

export type UseUnreadsHookReturn = {
  unreads: string[];
  setUnreads: (unreads: string[]) => void;
  addUnreadMessageIds: (unreadMsgIds: string[]) => void;
  updateUnreadMessageIds: (currentlyVisibleMessages: string[]) => void;
  getUnreadMessages: () => Set<string>;
};

function doLoad(): string[] {
  if (typeof window === "undefined") {
    return [];
  }
  const stored = localStorage.getItem(UNREADS_STORAGE_KEY);
  if (!stored) {
    return [];
  }

  try {
    return stored.split(",") as string[];
  } catch {
    return [];
  }
}

// this hook maintains a globally shared pool of unread message IDs.
export function useUnreads(): UseUnreadsHookReturn {
  const [unreads, setUnreads] = useState<string[] | undefined>(undefined);

  function doStore(unreadMsgIds: string[] | Set<string>) {
    const ids = Array.from(unreadMsgIds);
    localStorage.setItem(UNREADS_STORAGE_KEY, ids.join(","));
    setUnreads(ids);
  }

  const getUnreadMessages = (): Set<string> => {
    return new Set(doLoad());
  };

  const addUnreadMessageIds = (unreadMsgIds: string[]) => {
    if (typeof window === "undefined") {
      return;
    }

    const newUnreads = Array.from(
      new Set([...getUnreadMessages(), ...unreadMsgIds]),
    );
    doStore(newUnreads);
  };

  const updateUnreadMessageIds = (currentlyVisibleMessages: string[]) => {
    if (typeof window === "undefined") {
      return new Set();
    }
    const existing = getUnreadMessages();
    const visibleSet = new Set(currentlyVisibleMessages);
    const remaining = existing.difference(visibleSet);
    doStore(remaining);
  };

  return {
    unreads: unreads || doLoad(),
    setUnreads: (unreads) => doStore(unreads),
    addUnreadMessageIds,
    updateUnreadMessageIds,
    getUnreadMessages,
  };
}
