import type { Metadata } from "next";
import { cookies } from 'next/headers';
import "./globals.css";
import Header from "@/components/Header";
import { identityAPI } from "@/lib/api/client";

async function getCurrentUser() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session')?.value;

    if (!token) {
        return null;
    }

    try {
        const data = await identityAPI.getProfile(token);
        return data.user;
    } catch (error) {
        console.error("Failed to fetch current user in layout:", error);
        return null;
    }
}

export const metadata: Metadata = {
    title: "فروشگاه بلک شاپ",
    description: "یک فروشگاه مدرن ساخته شده با Next.js و Go",
};


export default async function RootLayout({
                                             children,
                                         }: Readonly<{
    children: React.ReactNode;
}>) {
    const user = await getCurrentUser();

    return (
        <html lang="fa" dir="rtl">
        <body className="bg-gray-50 text-gray-800">
        <Header user={user} />
        <main className="container mx-auto p-4 sm:p-6 lg:p-8">
            {children}
        </main>
        </body>
        </html>
    );
}
