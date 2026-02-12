"use client";

import { ConnEntry } from "@/apis/types";
import { MenuItem, Typography } from "@mui/material";

export function RenderPeerEntry(props: {
  conn: ConnEntry;
  activeNodeId: string;
  onSelect: () => void;
}) {
  const { conn, activeNodeId, onSelect } = props;
  return (
    <MenuItem
      selected={activeNodeId === conn.node_id}
      onClick={() => {
        onSelect();
      }}
      sx={{ overflow: "hidden" }}
    >
      {conn.entry?.node_name || conn.node_id}
      <Typography
        component="span"
        variant="body2"
        gutterBottom={false}
        marginLeft={1}
        noWrap
      >
        {conn.node_id}
      </Typography>
    </MenuItem>
  );
}
