import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { Vazirmatn } from "next/font/google";
import { cn } from "@/lib/utils";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const fontSans = Vazirmatn({
    subsets: ["latin", "arabic"],
    variable: "--font-sans",
});

export const metadata: Metadata = {
  title: "پنل مدیریت فروشگاه",
  description: "پنل مدیریت فروشگاه بلک شاپ",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="fa" dir="rtl">
      <body
          className={cn(
              "min-h-screen bg-background font-sans antialiased",
              fontSans.variable
          )}
      >
      {children}
      </body>
    </html>
  );
}
