import { Dispatch, SetStateAction } from "react";
import { PSKey, usePersistentStorage } from "./persistent";
import { Preference } from "./types";

export function usePreference() {
  const prefSt = usePersistentStorage(PSKey.Preference);
  const getVal: () => Preference = () => {
    try {
      const x = JSON.parse(prefSt.getValue() || "{}");
      if (typeof x === "object") {
        return x;
      }
    } catch (_) {}
    return {};
  };

  const setPreference: Dispatch<SetStateAction<Preference>> = (updater) => {
    if (typeof updater === "function" || updater instanceof Function) {
      const newVal = updater(getVal());
      prefSt.setValue(JSON.stringify(newVal));
    } else {
      prefSt.setValue(JSON.stringify(updater));
    }
  };

  return {
    preference: getVal(),
    setPreference,
  };
}
