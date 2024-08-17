import type { Metadata } from "next";
import NavBar from "@/app/ui/Nav/NavBar";
import "./globals.css";

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
      <body>
        <main className="min-h-dvh flex flex-col items-center">
          <NavBar />
          {children}
        </main>
      </body>
    </html>
  );
}
