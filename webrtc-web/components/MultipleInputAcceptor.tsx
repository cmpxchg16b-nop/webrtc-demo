"use client";

import {
  Button,
  Dialog,
  DialogActions,
  DialogTitle,
  DialogContent,
  TextField,
} from "@mui/material";
import { useState } from "react";

export function MultipleInputAcceptor(props: {
  title: string;
  rows?: number;
  open: boolean;
  onCancel: () => void;
  onConfirm: (input: string) => void;
}) {
  const [candidateText, setCandidateText] = useState("");
  const { title, rows = 4, open, onCancel, onConfirm } = props;
  return (
    <Dialog
      maxWidth="md"
      fullWidth
      open={open}
      onClose={() => {
        onCancel();
      }}
    >
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <TextField
          variant="outlined"
          multiline
          rows={rows}
          fullWidth
          value={candidateText}
          onChange={(e) => {
            setCandidateText(e.target.value);
          }}
        />
        <DialogActions sx={{ marginTop: 2 }}>
          <Button
            onClick={() => {
              onCancel();
            }}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={() => {
              onConfirm(candidateText);
            }}
          >
            Set
          </Button>
        </DialogActions>
      </DialogContent>
    </Dialog>
  );
}
