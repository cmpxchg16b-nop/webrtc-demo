"use client";

import { Box, Tooltip } from "@mui/material";
import { getPreferredColor } from "./ChangePreference";
import { useQuery } from "@tanstack/react-query";
import { getAvatar } from "@/apis/profile";
import {
  getColorTokenHashFromUsername,
  paintFirstLetterAvatar,
  PRESET_COLORS,
} from "@/apis/colors";
import { AuthenticationType } from "@/apis/types";

export function RenderAvatar(props: {
  username: string;
  size?: "default" | "small" | "large";
  preferredColorIdx?: number | string;
  authentication: AuthenticationType | undefined;
}) {
  const { username, size = "default", authentication } = props;
  const firstCap =
    username && username.length > 0 ? username[0].toUpperCase() : "";

  const variants = {
    large: "64px",
    default: "48px",
    small: "32px",
  };

  const fontSizeVariants = {
    large: "2rem",
    default: "1.5rem",
    small: "1rem",
  };

  let bgColorUsedLight: string = "orange";
  let bgColorUsedDark: string = "orange";
  const preferredColorIdx =
    props.preferredColorIdx === undefined || props.preferredColorIdx === null
      ? getColorTokenHashFromUsername(username, PRESET_COLORS.length)
      : props.preferredColorIdx;

  const colorToken = getPreferredColor(preferredColorIdx);
  bgColorUsedLight = colorToken.light;
  bgColorUsedDark = colorToken.dark;

  // Use React Query to fetch avatar from IAPOperator
  const { data: avatarUrl } = useQuery({
    queryKey: ["avatar", username],
    queryFn: async () => {
      try {
        const dataUrl = await getAvatar(username);
        return dataUrl;
      } catch (error) {
        console.error("Failed to fetch avatar from IAPOperator:", error);
        return paintFirstLetterAvatar(username || "");
      }
    },
  });

  const normalAvatar = (
    <Box
      component="img"
      src={avatarUrl}
      alt={username}
      sx={{
        width: variants[size],
        height: variants[size],
        borderRadius: "100%",
        objectFit: "cover",
        flexShrink: 0,
      }}
    />
  );

  const fallBackAvatar = (
    <Box
      sx={[
        {
          width: variants[size],
          height: variants[size],
          backgroundColor: bgColorUsedLight,
          borderRadius: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontWeight: "bold",
          fontSize: fontSizeVariants[size],
          flexShrink: 0,
          color: "white",
        },
        (theme) =>
          theme.applyStyles("dark", {
            backgroundColor: bgColorUsedDark,
          }),
      ]}
    >
      {firstCap}
    </Box>
  );

  return (
    <Tooltip
      title={
        <Box>
          <Box>{`Username: @${username}`}</Box>
          <Box>{`Authentication: ${authentication}`}</Box>
        </Box>
      }
    >
      {avatarUrl ? normalAvatar : fallBackAvatar}
    </Tooltip>
  );
}
