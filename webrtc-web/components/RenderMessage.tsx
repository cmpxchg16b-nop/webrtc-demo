"use client";

import { ChatMessage } from "@/apis/types";
import { Box, Card } from "@mui/material";

export function RenderMessage(props: { message: ChatMessage }) {
  const { message } = props;
  return (
    <Card
      sx={{
        padding: 2,
        gap: 1,
        flexWrap: "wrap",
        justifyContent: "space-between",
        alignItems: "center",
        maxWidth: "100%",
        width: "max-content",
        flexShrink: 0,
      }}
    >
      <Box sx={{ whiteSpace: "pre-wrap" }}>{message.message}</Box>
    </Card>
  );
}
