"use client";

import {
  ChatMessage,
  ChatMessageFile,
  FileTransferStatusEntry,
} from "@/apis/types";
import { InsertDriveFile } from "@mui/icons-material";
import { Box, Card } from "@mui/material";
import { Fragment } from "react/jsx-runtime";

function getFileLoadedRatio(
  file: ChatMessageFile,
  fileTransferStatus: Record<string, FileTransferStatusEntry>,
): number | undefined {
  if (file.dcId) {
    const fileLoadingStatus = fileTransferStatus[file.dcId];
    if (fileLoadingStatus) {
      const fileTotalSize = Math.max(1, file.size ?? 0);
      return Math.min(1, fileLoadingStatus.bytesReceived / fileTotalSize);
    }
  }
}

function RenderFile(props: {
  file: ChatMessageFile;
  fileTransferStatus: Record<string, FileTransferStatusEntry>;
}) {
  const { file, fileTransferStatus } = props;
  const fileLoadingStatus = file.dcId
    ? fileTransferStatus[file.dcId]
    : undefined;
  const loadedRatio: number | undefined = getFileLoadedRatio(
    file,
    fileTransferStatus,
  );

  return (
    <Fragment>
      {loadedRatio !== undefined &&
        loadedRatio !== null &&
        fileLoadingStatus &&
        !fileLoadingStatus?.closed && (
          <Box
            sx={{
              position: "absolute",
              top: 0,
              right: 0,
              width: `${(1 - loadedRatio) * 100}%`,
              height: "100%",
              backgroundColor: "rgba(0, 0, 0, 0.5)",
            }}
          ></Box>
        )}
      <Box sx={{ padding: 2 }}>
        {file.url ? (
          <a href={file.url} download={file.name}>
            <InsertDriveFile />
            <Box component="span" sx={{ paddingLeft: 0.5 }}>
              {file.name}
            </Box>
          </a>
        ) : (
          (file.name ?? "(unknown file)")
        )}{" "}
        {loadedRatio !== undefined &&
          loadedRatio !== null &&
          fileLoadingStatus &&
          !fileLoadingStatus?.closed && (
            <Box component="span" sx={{ paddingLeft: 0.5 }}>
              {`(${Math.round(loadedRatio * 100)}%)`}
            </Box>
          )}
      </Box>
    </Fragment>
  );
}

export function RenderMessage(props: {
  message: ChatMessage;
  onAmend?: (amendedMsg: ChatMessage) => void;
  onDelete?: (deletedMsgId: string) => void;
  fileTransferStatus: Record<string, FileTransferStatusEntry>;
}) {
  // todo: add message edit feature and delete feature in context menu
  const { message, onAmend, onDelete, fileTransferStatus } = props;

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
          <RenderFile
            file={message.file}
            fileTransferStatus={fileTransferStatus}
          />
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
