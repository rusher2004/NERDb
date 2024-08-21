import "./globals.css";
import type { Metadata } from "next";
import NavBar from "@/app/ui/Nav/NavBar";
import { exo2 } from "@/app/ui/fonts";
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
