"use client";

import {
  ChatMessage,
  ChatMessageFile,
  ChatMessageFileCategory,
  ConnEntry,
  ConnTrackStatus,
  ConnTrackStatusEntry,
} from "@/apis/types";
import { Badge, Box, MenuItem, Typography } from "@mui/material";
import { RenderAvatar } from "./RenderAvatar";

function getMessagePreviewThumbnail(msgCat: ChatMessageFileCategory): string {
  if (msgCat === ChatMessageFileCategory.Image) {
    return "🖼";
  } else if (msgCat === ChatMessageFileCategory.Video) {
    return "🎥";
  } else if (msgCat === ChatMessageFileCategory.Audio) {
    return "🎵";
  } else if (msgCat === ChatMessageFileCategory.Document) {
    return "📄";
  } else {
    return "📎";
  }
}

function getMessagePreview(msg: ChatMessage): string {
  if (msg.message) {
    return msg.message.length > 30
      ? msg.message.slice(0, 30) + "..."
      : msg.message;
  }
  if (msg.file) {
    const thumbnail = getMessagePreviewThumbnail(msg.file.category);
    return msg.file.name
      ? `${thumbnail} ${msg.file.name}`
      : `${thumbnail} File`;
  }
  return "";
}

export function RenderPeerEntry(props: {
  conn: ConnEntry;
  avatarUrl?: string;
  activeNodeId: string;
  onSelect: () => void;
  rtt?: number;
  latestUnreadMessage?: ChatMessage;
  numUnreads?: number;
}) {
  const {
    conn,
    avatarUrl,
    activeNodeId,
    onSelect,
    rtt,
    latestUnreadMessage,
    numUnreads,
  } = props;
  const hasUnreads = numUnreads !== undefined && numUnreads > 0;

  return (
    <MenuItem
      selected={activeNodeId === conn.node_id}
      onClick={() => {
        onSelect();
      }}
      sx={{
        overflow: "hidden",
        display: "flex",
        alignItems: "flex-start",
        gap: 1,
        flexDirection: "column",
      }}
    >
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          gap: 1,
          width: "100%",
        }}
      >
        <Badge
          badgeContent={numUnreads}
          color="primary"
          invisible={!hasUnreads}
          max={99}
        >
          <RenderAvatar
            username={conn.entry?.node_name || conn.node_id}
            url={avatarUrl}
            size="small"
          />
        </Badge>
        <Box sx={{ flex: 1, minWidth: 0 }}>
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            <Typography
              noWrap
              sx={{ fontWeight: hasUnreads ? "bold" : "normal" }}
            >
              {conn.entry?.node_name || conn.node_id}
            </Typography>
            {rtt !== undefined && (
              <Typography
                component="span"
                variant="body2"
                color="text.secondary"
                noWrap
              >
                {rtt.toFixed(2).replace(/\.?0+$/, "")}ms
              </Typography>
            )}
          </Box>
          {hasUnreads && latestUnreadMessage && (
            <Typography
              variant="body2"
              color="text.secondary"
              noWrap
              sx={{ fontWeight: hasUnreads ? "medium" : "normal" }}
            >
              {getMessagePreview(latestUnreadMessage)}
            </Typography>
          )}
        </Box>
      </Box>
    </MenuItem>
  );
}
