"use client";

import { ChatMessage } from "@/apis/types";
import { Box, IconButton, Input, Tooltip } from "@mui/material";
import { useState } from "react";
import SendIcon from "@mui/icons-material/Send";

export function MessageComposer(props: {
  onMessage: (message: ChatMessage) => void;
}) {
  const [messageInput, setMessageInput] = useState<string>("");
  const doSend = () => {
    const msgTxt = messageInput;

    const msgObject: ChatMessage = {
      messageId: crypto.randomUUID(),
      fromNodeId: "",
      toNodeId: "",
      message: msgTxt,
      timestamp: Date.now(),
    };

    props.onMessage(msgObject);

    setMessageInput("");
  };

  const [shiftPressed, setShiftPressed] = useState(false);

  const handleKeyUp = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Shift") {
      setShiftPressed(false);
    }
  };

  const handleEnterKeyPress = (
    event: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === "Shift") {
      setShiftPressed(true);
      return;
    }
    if (event.key === "Enter" && !shiftPressed) {
      event.preventDefault();
      event.stopPropagation();
      doSend();
    }
  };

  return (
    <Box
      sx={{
        borderTop: "1px solid #999",
        display: "flex",
        flexDirection: "row",
        alignItems: "flex-end",
        paddingTop: 1,
        paddingBottom: 1,
        paddingLeft: 1.5,
        paddingRight: 1,
      }}
    >
      <Input
        fullWidth
        multiline
        maxRows={8}
        value={messageInput}
        onChange={(e) => setMessageInput(e.target.value)}
        onKeyDown={handleEnterKeyPress}
        onKeyUp={handleKeyUp}
        disableUnderline
      />
      <Tooltip title="Send">
        <IconButton
          onClick={() => {
            doSend();
          }}
        >
          <SendIcon />
        </IconButton>
      </Tooltip>
    </Box>
  );
}
