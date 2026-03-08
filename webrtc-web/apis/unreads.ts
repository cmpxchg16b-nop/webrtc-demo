import { RefObject, useState } from "react";
import { PSKey, usePersistentStorage } from "./persistent";

const UNREADS_STORAGE_KEY = "webrtc_unread_message_ids";

export type UseUnreadsHookReturn = {
  // map of <to_node_id>, <msg_ids>
  unreads: Record<string, string[]>;
  setUnreads: (to_node_id: string, unreads: string[]) => void;
  addUnreadMessageIds: (to_node_id: string, unreadMsgIds: string[]) => void;
  updateUnreadMessageIds: (
    to_node_id: string,
    currentlyVisibleMessages: string[],
  ) => void;
  getUnreadMessages: () => Record<string, string[]>;
};

// this hook maintains a globally shared pool of unread message IDs.
// and it serves as the single authority of unread message IDs,
// any message isn't really unread unless it has been queued into here,
// any message isn't really read unless it has been removed from here.
export function useUnreads(): UseUnreadsHookReturn {
  const unreadsSt = usePersistentStorage(PSKey.Unreads);
  const getVal: () => Record<string, string[]> = () => {
    const storedVal = unreadsSt.getValue();
    if (storedVal) {
      try {
        return (JSON.parse(storedVal) || {}) as Record<string, string[]>;
      } catch (_) {}
    }
    return {};
  };

  return {
    unreads: getVal(),
    setUnreads: (toNode, unreads) => {
      const val = getVal();
      const newVal = {
        ...val,
        [toNode]: unreads,
      };
      unreadsSt.setValue(JSON.stringify(newVal));
    },
    addUnreadMessageIds: (toNode, unreads) => {
      if (!toNode) {
        console.error("add unread msgids to null node");
        return;
      }
      const val = getVal();
      const newVal = {
        ...val,
        [toNode]: [...(val[toNode] ?? []), ...unreads],
      };
      unreadsSt.setValue(JSON.stringify(newVal));
    },
    updateUnreadMessageIds: (toNode, visibles) => {
      const val = getVal();
      const originSet = new Set(val[toNode] ?? []);
      const visibleSet = new Set(visibles);
      const diffSet = originSet.difference(visibleSet);
      const newVal = {
        ...val,
        [toNode]: Array.from(diffSet),
      };
      unreadsSt.setValue(JSON.stringify(newVal));
    },
    getUnreadMessages: () => {
      return getVal();
    },
  };
}
