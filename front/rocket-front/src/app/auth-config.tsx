// auth-config.tsx
import { UserManager, UserManagerSettings, WebStorageStateStore } from "oidc-client-ts";

let userManager: UserManager | null = null;

const getBaseUrl = () => {
    if (typeof window !== 'undefined') {
        return `${window.location.protocol}//${window.location.host}`;
    }
    return "http://localhost:3000"; // Fallback for SSR
};

const getUserManager = (): UserManager => {
    if (typeof window !== 'undefined' && !userManager) {
        const oidcConfig: UserManagerSettings = {
            authority: "https://login.microsoftonline.com/af22c6ec-bb59-46c1-b14e-b0335ee947a2/v2.0",
            client_id: "a487944c-e6e5-472b-a394-276fffda940a",
            redirect_uri: `${getBaseUrl()}/callback`,
            post_logout_redirect_uri: "http://localhost:3000",
            response_type: "code",
            scope: "openid profile email",
            loadUserInfo: true,
            userStore: new WebStorageStateStore({ store: window.localStorage })
        };

        userManager = new UserManager(oidcConfig);
    }

    return userManager as UserManager;
};

export default getUserManager;