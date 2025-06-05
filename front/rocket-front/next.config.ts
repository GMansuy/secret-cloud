import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    output: 'export',
    // Needed for static export
    images: {
        unoptimized: true,
    },
};

export default nextConfig;
