"use client";

import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
} from "@mui/material";
import CheckIcon from "@mui/icons-material/Check";
import { Dispatch, SetStateAction, useState } from "react";

export const PRESET_COLORS = [
  "#F44336", // Red
  "#E91E63", // Pink
  "#9C27B0", // Purple
  "#673AB7", // Deep Purple
  "#3F51B5", // Indigo
  "#2196F3", // Blue
  "#00BCD4", // Cyan
  "#009688", // Teal
  "#4CAF50", // Green
  "#FF9800", // Orange
];

export type Preference = {
  name: string;
  indexOfPreferColor: number;
};

export function ChangePreference(props: {
  value: Preference;
  onChange: Dispatch<SetStateAction<Preference>>;
  open: boolean;
  onClose: () => void;
  onConfirm: (name: string) => Promise<void>;
}) {
  const { value, onChange, open, onClose, onConfirm } = props;
  const { name, indexOfPreferColor } = value;
  const onNameChange = (name: string) => {
    onChange({ ...value, name });
  };
  const onColorChange = (index: number) => {
    onChange({ ...value, indexOfPreferColor: index });
  };
  const [waiting, setWaiting] = useState(false);
  return (
    <Dialog maxWidth="sm" fullWidth open={open} onClose={onClose}>
      <DialogTitle>Preference</DialogTitle>
      <DialogContent>
        <TextField
          variant="standard"
          label="Name"
          fullWidth
          value={name}
          onChange={(e) => {
            onNameChange(e.target.value);
          }}
        />
        <Box sx={{ marginTop: 2 }}>
          <Typography
            variant="subtitle2"
            color="text.secondary"
            sx={{ marginBottom: 1 }}
          >
            Color
          </Typography>
          <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1.5 }}>
            {PRESET_COLORS.map((color, index) => (
              <Box
                key={index}
                onClick={() => onColorChange(index)}
                sx={{
                  width: 36,
                  height: 36,
                  borderRadius: "50%",
                  backgroundColor: color,
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  transition: "transform 0.1s ease-in-out",
                  border: indexOfPreferColor === index ? "2px solid" : "none",
                  borderColor:
                    indexOfPreferColor === index
                      ? "primary.main"
                      : "transparent",
                }}
              >
                {indexOfPreferColor === index && (
                  <CheckIcon sx={{ color: "white", fontSize: 20 }} />
                )}
              </Box>
            ))}
          </Box>
        </Box>
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
