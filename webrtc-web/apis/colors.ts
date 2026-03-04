export type ColorToken = {
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

export function getColorTokenHashFromUsername(username: string): number {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    const char = username.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash; // Convert to 32bit integer
  }
  return Math.abs(hash) % PRESET_COLORS.length;
}
