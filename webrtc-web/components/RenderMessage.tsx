"use client";

import { ChatMessage, FileTransferStatusEntry } from "@/apis/types";
import { InsertDriveFile } from "@mui/icons-material";
import { Box, Card } from "@mui/material";

export function RenderMessage(props: {
  message: ChatMessage;
  onAmend?: (amendedMsg: ChatMessage) => void;
  onDelete?: (deletedMsgId: string) => void;
  fileTransferStatus: Record<string, FileTransferStatusEntry>;
}) {
  // todo: add message edit feature and delete feature in context menu
  const { message, onAmend, onDelete, fileTransferStatus } = props;
  let loadingProgress = "";
  const dcId = message.file?.dcId;
  const fileLoadingStatus = dcId ? fileTransferStatus[dcId] : undefined;

  let fileLoadedRatio: number | undefined;
  if (fileLoadingStatus) {
    const fileTotalSize = message.file?.size ?? 0;
    fileLoadedRatio = fileLoadingStatus.bytesReceived / fileTotalSize;
    const percentage = Math.round(fileLoadedRatio * 100);
    loadingProgress = `(${percentage}%)`;
  }

  return (
    <Box>
      <Card
        sx={{
          gap: 1,
          flexWrap: "wrap",
          justifyContent: "space-between",
          alignItems: "center",
          maxWidth: "100%",
          width: "max-content",
          flexShrink: 0,
          position: "relative",
        }}
      >
        {fileLoadedRatio && fileLoadingStatus && !fileLoadingStatus?.closed && (
          <Box
            sx={{
              position: "absolute",
              top: 0,
              right: 0,
              width: `${(1 - fileLoadedRatio) * 100}%`,
              height: "100%",
              backgroundColor: "rgba(0, 0, 0, 0.5)",
            }}
          ></Box>
        )}
        {message.image && (
          <img
            style={{ maxHeight: "240px" }}
            src={message.image.url}
            alt={message.message}
          />
        )}
        {message.video && (
          <video
            autoPlay={false}
            controls
            style={{ maxHeight: "240px" }}
            src={message.video.url}
          />
        )}
        {message.file && (
          <Box sx={{ padding: 2 }}>
            <a href={message.file.url} download={message.file.name}>
              <InsertDriveFile />
              <Box component="span" sx={{ paddingLeft: 0.5 }}>
                {message.file.name}
              </Box>
              {loadingProgress &&
                fileLoadingStatus &&
                !fileLoadingStatus?.closed && (
                  <Box component="span" sx={{ paddingLeft: 0.5 }}>
                    {loadingProgress}
                  </Box>
                )}
            </a>
          </Box>
        )}
        {message.message && (
          <Box sx={{ padding: 2, whiteSpace: "pre-wrap" }}>
            {message.message}
          </Box>
        )}
        {message.richText && (
          <Box sx={{ padding: 2, whiteSpace: "pre-wrap" }}>
            {message.richText.content}
          </Box>
        )}
      </Card>
    </Box>
  );
}
