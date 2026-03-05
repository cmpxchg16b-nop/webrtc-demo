"use client";

import { WSServer } from "@/apis/types";
import { Box, TextField, Select, MenuItem, Button } from "@mui/material";
import { IaPLoginButton } from "./LoginButton";

// Select what signalling server to use

export function ServerSelector(props: {
  servers: WSServer[];
  onConnect: (server: WSServer) => void;
  selectedServer: string;
  onSelectedServerChange: (serverId: string) => void;
  preferName: string;
  onPreferNameChange: (preferName: string) => void;
  connecting: boolean;
}) {
  const {
    servers,
    selectedServer,
    onSelectedServerChange,
    onConnect,
    preferName,
    onPreferNameChange,
    connecting,
  } = props;

  const selectedServerObj = servers.find(
    (server) => server.id === selectedServer,
  );
  const hasIAP = selectedServerObj?.iap && selectedServerObj.iap.loginUrl;

  const handleLoginClick = () => {
    if (hasIAP && selectedServerObj?.iap?.loginUrl) {
      window.location.href = selectedServerObj.iap.loginUrl;
    }
  };

  const getDisplayName = () => {
    if (!selectedServerObj?.iap) return "Connect";
    const displayName = selectedServerObj.iap.displayName;
    if (typeof displayName === "string") {
      return displayName;
    }
    return displayName.en_US;
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        height: "100%",
      }}
    >
      <Box
        sx={{
          display: "grid",
          gridTemplateColumns: "auto 1fr",
          gap: 1,
          rowGap: 2,
          alignItems: "center",
          padding: 2,
        }}
      >
        <Box sx={{ justifySelf: "right" }}>Choose Server:</Box>
        <Select
          variant="standard"
          label="Server"
          value={selectedServer}
          onChange={(e) => onSelectedServerChange(e.target.value)}
        >
          {servers.map((server) => (
            <MenuItem key={server.id} value={server.id}>
              {server.name}
            </MenuItem>
          ))}
        </Select>
        {!hasIAP && (
          <>
            <Box sx={{ justifySelf: "right" }}>Pick a Name:</Box>
            <TextField
              fullWidth
              variant="standard"
              value={preferName}
              onChange={(e) => onPreferNameChange(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  e.stopPropagation();
                  const server = servers.find(
                    (server) => server.id === selectedServer,
                  );
                  if (server) {
                    onConnect(server);
                  }
                }
              }}
            />
          </>
        )}
      </Box>
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          marginTop: 2,
        }}
      >
        {hasIAP ? (
          <IaPLoginButton
            onClick={handleLoginClick}
            iapContext={selectedServerObj!.iap!}
            loading={connecting}
          />
        ) : (
          <Button
            variant="contained"
            loading={connecting}
            onClick={() => {
              const server = servers.find(
                (server) => server.id === selectedServer,
              );
              if (server) {
                onConnect(server);
              }
            }}
          >
            Connect
          </Button>
        )}
      </Box>
    </Box>
  );
}
