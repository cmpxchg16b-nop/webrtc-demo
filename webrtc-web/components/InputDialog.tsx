"use client";

import { MultipleInputAcceptor } from "@/components/MultipleInputAcceptor";
import { ContentCopy, Refresh } from "@mui/icons-material";
import {
  Box,
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  Tooltip,
  IconButton,
  DialogActions,
  Button,
} from "@mui/material";
import { useState, RefObject } from "react";

export function RemoteDescriptionInputDialog(props: {
  peerConnectionRef: RefObject<RTCPeerConnection | null>;
  open: boolean;
  onClose: () => void;
}) {
  const { peerConnectionRef, open, onClose } = props;
  return (
    <MultipleInputAcceptor
      title="Remote Description"
      rows={4}
      open={open}
      onCancel={onClose}
      onConfirm={(input) => {
        try {
          const remoteDescription = JSON.parse(input);
          peerConnectionRef.current?.setRemoteDescription(remoteDescription);
          onClose();
        } catch (e) {
          console.error(e);
        }
      }}
    />
  );
}

export function LocalDescriptionInputDialog(props: {
  peerConnectionRef: RefObject<RTCPeerConnection | null>;
  open: boolean;
  onClose: () => void;
}) {
  const { peerConnectionRef, open, onClose } = props;
  return (
    <MultipleInputAcceptor
      title="Local Description"
      rows={4}
      open={open}
      onCancel={onClose}
      onConfirm={(input) => {
        try {
          const localDescription = JSON.parse(input);
          peerConnectionRef.current?.setLocalDescription(localDescription);
          onClose();
        } catch (e) {
          console.error(e);
        }
      }}
    />
  );
}

export function OfferDialog(props: {
  open: boolean;
  onClose: () => void;
  peerConnectionRef: RefObject<RTCPeerConnection | null>;
  dataChannelRef: RefObject<RTCDataChannel | null>;
}) {
  const [candidateText, setCandidateText] = useState("");
  const { open, onClose, peerConnectionRef, dataChannelRef } = props;
  return (
    <Dialog
      maxWidth="md"
      fullWidth
      open={open}
      onClose={() => {
        onClose();
      }}
    >
      <DialogTitle
        sx={{
          display: "flex",
          flexDirection: "row",
          gap: 1,
          alignItems: "center",
          justifyContent: "space-between",
          flexWrap: "wrap",
        }}
      >
        <Box>Offer</Box>
        <Box>
          <Tooltip title={"Refresh"}>
            <IconButton
              onClick={() => {
                const peerConnection = peerConnectionRef.current;
                const dc = peerConnection?.createDataChannel("dc1");
                if (dc) {
                  dataChannelRef.current = dc;
                  dc.onopen = () => {
                    console.log("[dbg] data channel opened", dc);
                  };
                  dc.onclose = () => {
                    console.log("[dbg] data channel closed", dc);
                  };
                  dc.onerror = (error) => {
                    console.error("[dbg] data channel error", error);
                  };
                  dc.onmessage = (event) => {
                    console.log("[dbg] data channel message", event.data, dc);
                  };

                  peerConnectionRef.current?.createOffer().then((offer) => {
                    setCandidateText(JSON.stringify(offer));
                    peerConnection?.setLocalDescription(offer);
                  });
                }
              }}
            >
              <Refresh />
            </IconButton>
          </Tooltip>
          <Tooltip title={"Copy"}>
            <IconButton
              onClick={() => {
                navigator?.clipboard?.writeText(candidateText);
              }}
            >
              <ContentCopy />
            </IconButton>
          </Tooltip>
        </Box>
      </DialogTitle>
      <DialogContent>
        <TextField
          variant="outlined"
          multiline
          rows={4}
          fullWidth
          value={candidateText}
          onChange={(e) => {
            setCandidateText(e.target.value);
          }}
        />
      </DialogContent>
    </Dialog>
  );
}

export function CandidateInputDialog(props: {
  peerConnectionRef: RefObject<RTCPeerConnection | null>;
  open: boolean;
  onClose: () => void;
}) {
  const { peerConnectionRef, open, onClose } = props;
  return (
    <MultipleInputAcceptor
      title="Add Candidate"
      rows={4}
      open={open}
      onCancel={onClose}
      onConfirm={(input) => {
        try {
          const candidate = JSON.parse(input);
          peerConnectionRef.current?.addIceCandidate(candidate);
          onClose();
        } catch (e) {
          console.error(e);
        }
      }}
    />
  );
}

export function AnswerDialog(props: {
  open: boolean;
  onClose: () => void;
  peerConnectionRef: React.RefObject<RTCPeerConnection | null>;
}) {
  const { open, onClose, peerConnectionRef } = props;
  const [answerText, setAnswerText] = useState("");
  return (
    <Dialog
      maxWidth="md"
      fullWidth
      open={open}
      onClose={() => {
        onClose();
      }}
    >
      <DialogTitle
        sx={{
          display: "flex",
          flexDirection: "row",
          gap: 1,
          alignItems: "center",
          justifyContent: "space-between",
          flexWrap: "wrap",
        }}
      >
        <Box>Answer</Box>
        <Box>
          <Tooltip title={"Refresh"}>
            <IconButton
              onClick={() => {
                peerConnectionRef.current?.createAnswer().then((answer) => {
                  setAnswerText(JSON.stringify(answer));
                });
              }}
            >
              <Refresh />
            </IconButton>
          </Tooltip>
          <Tooltip title={"Copy"}>
            <IconButton
              onClick={() => {
                navigator?.clipboard?.writeText(answerText);
              }}
            >
              <ContentCopy />
            </IconButton>
          </Tooltip>
        </Box>
      </DialogTitle>
      <DialogContent>
        <TextField
          variant="outlined"
          multiline
          rows={4}
          fullWidth
          value={answerText}
        />
      </DialogContent>
    </Dialog>
  );
}

export function ChangeNameDialog(props: {
  name: string;
  onNameChange: (name: string) => void;
  open: boolean;
  onClose: () => void;
  onConfirm: (name: string) => Promise<void>;
}) {
  const { name, onNameChange, open, onClose, onConfirm } = props;
  const [waiting, setWaiting] = useState(false);
  return (
    <Dialog maxWidth="sm" fullWidth open={open} onClose={onClose}>
      <DialogTitle>Change Name</DialogTitle>
      <DialogContent>
        <TextField
          variant="standard"
          label="New Name"
          fullWidth
          value={name}
          onChange={(e) => {
            onNameChange(e.target.value);
          }}
        />
        <DialogActions sx={{ marginTop: 2 }}>
          <Button
            onClick={() => {
              onClose();
            }}
          >
            Cancel
          </Button>
          <Button
            loading={waiting}
            variant="contained"
            onClick={() => {
              setWaiting(true);
              onConfirm(name).finally(() => {
                setWaiting(false);
              });
            }}
          >
            Confirm
          </Button>
        </DialogActions>
      </DialogContent>
    </Dialog>
  );
}
