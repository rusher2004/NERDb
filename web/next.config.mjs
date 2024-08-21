/** @type {import('next').NextConfig} */

import { withSentryConfig } from '@sentry/nextjs';

const nextConfig = {};

export default withSentryConfig(nextConfig, {
  authToken: process.env.SENTRY_AUTH_TOKEN,
  org: "procyon-innovations",
  project: "nerdb-web",
});
