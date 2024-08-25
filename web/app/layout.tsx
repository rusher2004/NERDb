import "./globals.css";
import Head from "next/head";
import type { Metadata } from "next";
import { exo2 } from "@/app/ui/fonts";
import NavBar from "@/app/ui/Nav/NavBar";
import { Analytics } from "@vercel/analytics/react";
import { SpeedInsights } from "@vercel/speed-insights/next";

export const metadata: Metadata = {
  title: "New Eden Rivalry Database",
  description: "NERDb: An Eve Online bragging rights database",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <Head>
        <script
          defer
          src="https://cloud.umami.is/script.js"
          data-website-id="c94d6ed7-21a9-40aa-9cca-4c144b826d3c"
          key="umami-script"
        ></script>
      </Head>
      <body className={`${exo2.className} antialiased`}>
        <main className="min-h-dvh flex flex-col items-center gap-4">
          <NavBar />
          {children}
          <Analytics />
          <SpeedInsights />
        </main>
      </body>
    </html>
  );
}
