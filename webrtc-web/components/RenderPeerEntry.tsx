"use client";

import { ConnEntry, ConnTrackStatus, ConnTrackStatusEntry } from "@/apis/types";
import { MenuItem, Typography } from "@mui/material";
import { RenderAvatar } from "./RenderAvatar";

export function RenderPeerEntry(props: {
  conn: ConnEntry;
  avatarUrl?: string;
  activeNodeId: string;
  onSelect: () => void;
  rtt?: number;
}) {
  const { conn, avatarUrl, activeNodeId, onSelect, rtt } = props;
  return (
    <MenuItem
      selected={activeNodeId === conn.node_id}
      onClick={() => {
        onSelect();
      }}
      sx={{
        overflow: "hidden",
        display: "flex",
        alignItems: "center",
        gap: 1,
      }}
    >
      <RenderAvatar
        username={conn.entry?.node_name || conn.node_id}
        url={avatarUrl}
        size="small"
      />
      {conn.entry?.node_name || conn.node_id}
      {rtt !== undefined && (
        <Typography
          component="span"
          variant="body2"
          gutterBottom={false}
          marginLeft={1}
          noWrap
        >
          {rtt.toFixed(2).replace(/\.?0+$/, "")}ms
        </Typography>
      )}
    </MenuItem>
  );
}
