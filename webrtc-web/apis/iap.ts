import { DataURL, IAPKind } from "./types";

export interface IAPOperator {
  getAvatar(username: string): Promise<DataURL>;
}

export function mockIAPOperator() {
  return {
    getAvatar(username: string): Promise<DataURL> {
      // todo: mock it here
      return Promise.reject("unimplemented yet");
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
