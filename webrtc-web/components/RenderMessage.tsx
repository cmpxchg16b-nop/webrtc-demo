"use client";

import { ChatMessage } from "@/apis/types";
import { Box, Card } from "@mui/material";

export function RenderMessage(props: { message: ChatMessage }) {
  const { message } = props;
  return (
    <Card
      sx={{
        padding: 2,
        display: "flex",
        flexDirection: "row",
        gap: 1,
        flexWrap: "wrap",
        justifyContent: "space-between",
        alignItems: "center",
      }}
    >
      <Box>
        <Box>Message: {message.message}</Box>
      </Box>
    </Card>
  );
}
