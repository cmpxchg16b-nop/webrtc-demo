"use client";

import { WSServer } from "@/apis/types";
import {
  Box,
  TextField,
  Select,
  MenuItem,
  Button,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { IaPLoginButton } from "./LoginButton";
import { Fragment, useEffect } from "react";
import { PSKey, usePersistentStorage } from "@/apis/persistent";
import { useLoginStatusPolling } from "@/apis/profile";

const getNum = (s: string): number | undefined => {
  try {
    const x = parseInt(s);
    if (!Number.isNaN(x) && Number.isFinite(x)) {
      return x;
    }
  } catch (_) {}
};

const loginTimeoutMs = 60 * 1000;

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

  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const { getValue: getLoggingIn, setValue: setLoggingIn } =
    usePersistentStorage(PSKey.LoggingIn);
  const isLoggingIn = getLoggingIn() === "true";

  const { getValue: getLogInStart, setValue: setLogInStart } =
    usePersistentStorage(PSKey.LoggingStartedAt);

  const logInStartedAt = getLogInStart() || "";

  useEffect(() => {
    const loginStartTx = getNum(logInStartedAt);
    if (loginStartTx === undefined || loginStartTx === null) {
      return;
    }
    const loginTimeoutAt = loginStartTx + loginTimeoutMs;
    const now = new Date().valueOf();
    const timeDelta = Math.max(loginTimeoutAt - now, 0);
    const timeout = setTimeout(() => {
      setLoggingIn("false");
      setLogInStart("");
    }, timeDelta);
    return () => {
      clearTimeout(timeout);
    };
  }, [logInStartedAt]);

  const handleLoginClick = () => {
    // start polling (also the polling state would also survives page reload)
    setLoggingIn("true");
    setLogInStart(new Date().valueOf().toString());

    // navigate the user to the oauth2 authorization portal
    if (hasIAP && selectedServerObj?.iap?.loginUrl) {
      window.location.href = selectedServerObj.iap.loginUrl;
    }
  };

  const selectedServerObj = servers.find(
    (server) => server.id === selectedServer,
  );
  const hasIAP = selectedServerObj?.iap && selectedServerObj.iap.loginUrl;
  const handleConnect = () => {
    const server = servers.find((server) => server.id === selectedServer);
    if (server) {
      setLoggingIn("false");
      onConnect(server);
    }
  };

  const connectBtn = (
    <Button variant="contained" loading={connecting} onClick={handleConnect}>
      Connect
    </Button>
  );

  const { loggedIn, loggedInAs, hintText } = useLoginStatusPolling(
    selectedServerObj?.apiPrefix || "",
    3000,
  );

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
        <Box sx={{ justifySelf: "right" }}>{!isMobile && "Server"}</Box>
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
        {!hasIAP ? (
          <Fragment>
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
                  handleConnect();
                }
              }}
            />
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                marginTop: "2",
                gridColumn: "1 / span 2",
              }}
            >
              {connectBtn}
            </Box>
          </Fragment>
        ) : (
          <Fragment>
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                marginTop: "2",
                gridColumn: "1 / span 2",
                flexDirection: "column",
                gap: 2,
              }}
            >
              {isLoggingIn && (
                <Box
                  sx={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                >
                  {hintText}
                </Box>
              )}
              {loggedIn && loggedInAs ? (
                connectBtn
              ) : (
                <IaPLoginButton
                  loading={isLoggingIn}
                  onClick={handleLoginClick}
                  iapContext={selectedServerObj!.iap!}
                />
              )}
            </Box>
          </Fragment>
        )}
      </Box>
    </Box>
  );
}
