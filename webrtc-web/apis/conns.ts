import { ConnEntry, ConnRegistryData } from "./types";

const apiEndpoint = "http://localhost:3001";

export function getConns() {
  return fetch(`${apiEndpoint}/conns`)
    .then((r) => r.json())
    .then((r) => r as Record<string, ConnRegistryData>)
    .then((r) => {
      const entries: ConnEntry[] = [];
      return Object.entries(r).map(
        ([nodeId, entry]) =>
          ({
            node_id: nodeId,
            entry: entry,
          }) as ConnEntry,
      );
    });
}
