import { getProfile } from "@/apis/profile";
import { useQuery } from "@tanstack/react-query";

export function ShowDisplayName(props: {
  apiPrefix: string;
  username: string | undefined | null;
}) {
  const username = props.username ?? "";
  const { apiPrefix } = props;

  const { data } = useQuery({
    queryKey: ["profilefor", apiPrefix, "username", username],
    queryFn: () => getProfile(apiPrefix, username),
  });
  return <>{data?.displayName ?? username ?? ""}</>;
}
