import { paintFirstLetterAvatar, PRESET_COLORS } from "./colors";
import { DataURL, IAPKind } from "./types";

export interface IAPOperator {
  getAvatar(username: string): Promise<DataURL>;
}

function getDataURLFromBlob(blob: Blob): Promise<DataURL> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as DataURL);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}

export function mockIAPOperator(): IAPOperator {
  return {
    async getAvatar(username: string): Promise<DataURL> {
      const response = await fetch(
        `https://avatars.githubusercontent.com/${username}`,
      );

      try {
        const blob = await response.blob();
        const dataURL = await getDataURLFromBlob(blob);
        return dataURL;
      } catch (err) {
        console.error("failed to get avatar DataURL, falling back to default");
        return paintFirstLetterAvatar(username);
      }
    },
  };
}

export function getIAPOperator(kind: IAPKind): IAPOperator {
  switch (kind) {
    case IAPKind.MockIAP:
      return mockIAPOperator();
    default:
      throw new Error(`Unsupported IAP kind: ${kind}`);
  }
}
