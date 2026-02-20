import { RefObject, useCallback } from "react";

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
  const saveScrollTop = (
    peerId: string,
    elemRef: RefObject<HTMLDivElement | null>,
  ) => {
    if (elemRef.current) {
      const scrollTop = elemRef.current.scrollTop;
      const key = getScrollTopStorageKey(peerId);
      localStorage.setItem(key, String(scrollTop));
    }
  };

  const restoreScrollTop = (
    peerId: string,
    elemRef: RefObject<HTMLDivElement | null>,
  ) => {
    if (elemRef.current) {
      const key = getScrollTopStorageKey(peerId);
      const savedScrollTop = localStorage.getItem(key);
      if (savedScrollTop !== null) {
        elemRef.current.scrollTop = Number(savedScrollTop);
      }
    }
  };

  return {
    saveScrollTop,
    restoreScrollTop,
  };
}
