"use client";

import { Box, Button } from "@mui/material";

export function BasicWsInfo(props: {
  name?: string;
  url?: string;
  rtt?: number;
  nodeId?: string;
  lastSeq?: number;
  upTime?: number;
  onNameChangeRequested: () => void;
}) {
  const { name, url, rtt, nodeId, lastSeq, upTime, onNameChangeRequested } =
    props;

  return (
    <Box sx={{ padding: 2 }}>
      <Box>Basics Info</Box>
      <Box>
        Connected {url ? `to ${url}` : ""} {name ? `as ${name}` : ""}
      </Box>
      {nodeId && <Box>NodeId: {nodeId}</Box>}

      {rtt !== undefined && <Box>RTT: {rtt}ms</Box>}
      {lastSeq !== undefined && <Box>Last Seq: {lastSeq}</Box>}
      {upTime !== undefined && (
        <Box>
          Up Time:{" "}
          {(upTime / 1000).toFixed(3).replace(/0+$/, "").replace(/\.$/, "")}s
        </Box>
      )}

      <Box>
        <Button
          onClick={() => {
            onNameChangeRequested();
          }}
        >
          Change Name
        </Button>
      </Box>
    </Box>
  );
}
