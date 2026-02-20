import { RefObject } from "react";

export type UseScrollTopReturn = {
  // Save current scroll top of the specified element to the localStorage
  saveScrollTop: (
    peerId: string,
    elemRef: RefObject<HTMLDivElement | null>,
  ) => void;

  // Retrieve previously saved scroll top (if any) and apply it to the specified element
  restoreScrollTop: (
    peerId: string,
    elemRef: RefObject<HTMLDivElement | null>,
  ) => void;
};

export function getScrollTopStorageKey(peerId: string): string {
  return `saved_scroll_top:${peerId}`;
}

export function useScrollTop(): UseScrollTopReturn {
  // todo
}
