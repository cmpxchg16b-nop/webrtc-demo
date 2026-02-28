import { ConnEntry, ConnRegistryData } from "./types";

export function getConns(apiPrefix: string) {
  return fetch(`${apiPrefix}/conns`)
    .then((r) => r.json())
    .then((r) => r as Record<string, ConnRegistryData>)
    .then((r) => {
      const entries: ConnEntry[] = Object.entries(r).map(([nodeId, entry]) => {
        const registeredAt = entry.registered_at;
        return {
          node_id: nodeId,
          entry: entry,
          registered_at: registeredAt,
        } as ConnEntry;
      });
      entries.sort((a, b) => b.registered_at - a.registered_at);
      return entries;
    });
}
