"use client";

import { IDProvider, Preference, WSServer } from "@/apis/types";
import {
  Box,
  TextField,
  Select,
  MenuItem,
  Button,
  useMediaQuery,
  useTheme,
  Divider,
} from "@mui/material";
import { IdPLoginButton } from "./LoginButton";
import { Dispatch, Fragment, SetStateAction } from "react";
import {
  getLoginStatusHintTxt,
  getProfile,
  getProfileStatus,
} from "@/apis/profile";
import { useQuery } from "@tanstack/react-query";
import { PSKey, usePersistentStorage } from "@/apis/persistent";

// Select what signalling server to use
export function ServerSelector(props: {
  servers: WSServer[];
  onPinnedServerChange: (serverId: string) => void;
  connecting: boolean;
  onLogout: () => void;
  preference: Preference;
  onPreferenceChange: Dispatch<SetStateAction<Preference>>;
}) {
  const {
    servers,
    preference,
    onPreferenceChange,
    connecting,
    onPinnedServerChange,
    onLogout,
  } = props;

  const { getValue: getCurrentServer, setValue: setSelectedServer } =
    usePersistentStorage(PSKey.CurrentServer);

  // selectedServerId indicates the server that is currently active in the select box
  // the user might just selected a server, but didn't click the 'connect' button, so
  // the selectedServer might not necessarily be the pinnedSrv in the meantime
  const selectedServerId = getCurrentServer() || "";
  // pinnedSrv indicates which server the user decided to connect to

  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const selectedServerObj = servers.find(
    (server) => server.id === selectedServerId,
  );
  const handleLoginClick = (idp: IDProvider) => {
    // eslint-disable-next-line
    window.location.href = idp.loginUrl;
  };

  const handleConnect = () => {
    const server = servers.find((server) => server.id === selectedServerId);
    if (server) {
      // the app will automatically tries to connect to a pinned server
      onPinnedServerChange(server.id);
    }
  };

  const connectBtn = (
    <Button
      fullWidth
      variant="contained"
      loading={connecting}
      onClick={handleConnect}
    >
      Connect
    </Button>
  );
  const { isLoading: isLoginStatusLoading, data: profileStatusData } = useQuery(
    {
      queryKey: ["hasloggedin", selectedServerObj?.apiPrefix ?? ""],
      queryFn: () => getProfileStatus(selectedServerObj?.apiPrefix ?? ""),
    },
  );

  const { data: profileData } = useQuery({
    queryKey: ["profile", selectedServerObj?.apiPrefix ?? ""],
    queryFn: () => getProfile(selectedServerObj?.apiPrefix ?? ""),
  });

  const hintText = getLoginStatusHintTxt(
    profileStatusData?.logged_in,
    profileData,
  );
  const idps = selectedServerObj?.idp ?? [];

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
          value={selectedServerId}
          onChange={(e) => setSelectedServer(e.target.value)}
        >
          {servers.map((server) => (
            <MenuItem key={server.id} value={server.id}>
              {server.name}
            </MenuItem>
          ))}
        </Select>

        {isLoginStatusLoading ? (
          <Box>Fetching login status ...</Box>
        ) : profileStatusData?.logged_in ? (
          <Fragment>
            <Box>{hintText}</Box>
            <Box>{connectBtn}</Box>
            <Box>
              <Button fullWidth onClick={onLogout}>
                Logout
              </Button>
            </Box>
          </Fragment>
        ) : (
          <Fragment>
            {idps.map((idp) => (
              <IdPLoginButton
                key={idp.name}
                idpContext={idp}
                onClick={() => handleLoginClick(idp)}
              />
            ))}
          </Fragment>
        )}

        {selectedServerObj?.allowAnonymous && (
          <Box sx={{ paddingTop: 2 }}>
            <Divider orientation="horizontal" />
            <Box>Or, Connect as a visitor:</Box>
            <TextField
              fullWidth
              variant="standard"
              value={preference?.name ?? ""}
              onChange={(e) =>
                onPreferenceChange((prev) => {
                  return {
                    ...prev,
                    name: e.target.value,
                  };
                })
              }
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
          </Box>
        )}
      </Box>
    </Box>
  );
}
