"use client";

import { Box, Button, Card, Chip } from "@mui/material";

export function RenderICECandidate(props: { candidate: RTCIceCandidate }) {
  const { candidate } = props;
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
        <Box
          sx={{
            display: "flex",
            flexDirection: "row",
            gap: 1,
            flexWrap: "wrap",
          }}
        >
          <Chip label={`Type: ${candidate.type}`} />
          <Chip label={`Protocol: ${candidate.protocol}`} />
          <Chip label={`Component: ${candidate.component}`} />
          <Chip label={`Port: ${candidate.port}`} />
        </Box>
        <Box sx={{ marginTop: 1 }}>
          <Box>Address: {candidate.address}</Box>
        </Box>
      </Box>
      <Button
        onClick={() => {
          navigator?.clipboard?.writeText(JSON.stringify(candidate.toJSON()));
        }}
      >
        Copy
      </Button>
    </Card>
  );
}
