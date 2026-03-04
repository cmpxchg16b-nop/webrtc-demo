import { DataURL, IAPKind } from "./types";

export interface IAPOperator {
  getAvatar(username: string): Promise<DataURL>;
}

export function mockIAPOperator(): IAPOperator {
  return {
    async getAvatar(username: string): Promise<DataURL> {
      const response = await fetch(
        `https://avatars.githubusercontent.com/${username}`,
      );

      const blob = await response.blob();

      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result as DataURL);
        reader.onerror = reject;
        reader.readAsDataURL(blob);
      });
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
