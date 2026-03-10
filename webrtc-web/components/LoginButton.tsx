import { IAPKind, IDProvider } from "@/apis/types";
import { Button } from "@mui/material";
import { KioubitLogin } from "./web-components-declarative/KioubitLoginBtn";
import GithubIcon from "@mui/icons-material/github";

export function IdPLoginButton(props: {
  idpContext: IDProvider;
  onClick: () => void;
}) {
  const { idpContext, onClick } = props;

  switch (idpContext.kind) {
    case IAPKind.Kioubit:
      return <KioubitLogin onClick={onClick} />;
    case IAPKind.Github:
      return (
        <Button onClick={onClick} startIcon={<GithubIcon />}>
          Sign in with Github
        </Button>
      );
    default:
      return <Button onClick={onClick}>Login</Button>;
  }
}
