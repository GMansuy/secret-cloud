import { UserManager, UserManagerSettings,WebStorageStateStore } from "oidc-client-ts";

const oidcConfig: UserManagerSettings = {
    authority: "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/v2.0",
    client_id: "a487944c-e6e5-472b-a394-276fffda940a",
    redirect_uri: "http://localhost:3000/callback",
    post_logout_redirect_uri: "http://localhost:3000",
    response_type: "code",
    scope: "openid profile email", // add more if needed
    loadUserInfo: true,
    userStore: new WebStorageStateStore({ store: window.localStorage })
};

const userManager = new UserManager(oidcConfig);

export default userManager;