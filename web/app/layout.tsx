import "./globals.css";
import type { Metadata } from "next";
import NavBar from "@/app/ui/Nav/NavBar";
import { spaceMono } from "@/app/ui/fonts";
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
      <body className={`${spaceMono.className} antialiased`}>
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
