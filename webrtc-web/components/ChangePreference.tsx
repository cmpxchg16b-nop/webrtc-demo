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

type ColorToken = {
  light: string;
  dark: string;
};

export const PRESET_COLORS: ColorToken[] = [
  { light: "#EF9A9A", dark: "#F44336" }, // Red
  { light: "#F48FB1", dark: "#E91E63" }, // Pink
  { light: "#CE93D8", dark: "#9C27B0" }, // Purple
  { light: "#B39DDB", dark: "#673AB7" }, // Deep Purple
  { light: "#9FA8DA", dark: "#3F51B5" }, // Indigo
  { light: "#90CAF9", dark: "#2196F3" }, // Blue
  { light: "#80DEEA", dark: "#00BCD4" }, // Cyan
  { light: "#80CBC4", dark: "#009688" }, // Teal
  { light: "#A5D6A7", dark: "#4CAF50" }, // Green
  { light: "#FFCC80", dark: "#FF9800" }, // Orange
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
            {PRESET_COLORS.map((colorToken, index) => (
              <Box
                key={index}
                onClick={() => onColorChange(index)}
                sx={{
                  width: 36,
                  height: 36,
                  borderRadius: "50%",
                  cursor: "pointer",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  transition: "transform 0.1s ease-in-out",
                  transform: "rotate(45deg)",
                  background: `linear-gradient(to right, ${colorToken.light} 50%, ${colorToken.dark} 50%)`,
                  outline: indexOfPreferColor === index ? "3px solid" : "none",
                  outlineColor:
                    indexOfPreferColor === index
                      ? "primary.main"
                      : "transparent",
                }}
              ></Box>
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
