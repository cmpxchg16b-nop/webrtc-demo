export function logout(apiPrefix: string): Promise<void> {
  return fetch(`${apiPrefix}/logout`, {
    method: "POST",
    credentials: "include",
  }).then((r) => {
    console.log("Successfully logged out");
  });
}
