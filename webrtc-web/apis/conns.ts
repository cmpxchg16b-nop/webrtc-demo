import { ConnRegistryData } from "./types";

export function getConns(apiPrefix: string) {
  return fetch(`${apiPrefix}/conns`, { credentials: "include" })
    .then((r) => r.json())
    .then((r) => r as Record<string, ConnRegistryData>);
}
