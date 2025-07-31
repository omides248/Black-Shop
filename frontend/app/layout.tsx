// file: frontend/nextjs-app/app/layout.tsx
import type { Metadata } from "next";
import { cookies } from 'next/headers';
import "./globals.css";
import Header from "@/components/Header";
import { identityAPI, orderAPI } from "@/lib/api/client";

// // این تابع به صورت امن، اطلاعات اولیه کاربر و سبد خرید را از بک‌اند دریافت می‌کند
// async function getInitialData() {
//     // ✅ تغییر اصلی اینجاست: استفاده از await برای cookies()
//     const cookieStore = await cookies();
//     const token = cookieStore.get('session')?.value;
//
//     // اگر توکنی وجود نداشت، نیازی به ارسال درخواست نیست
//     if (!token) {
//         return { user: null, cart: null };
//     }
//
//     try {
//         // ما هر دو درخواست را به صورت همزمان ( موازی ) ارسال می‌کنیم تا سریع‌تر باشد
//         const [userResponse, cartResponse] = await Promise.all([
//             identityAPI.getProfile(token),
//             orderAPI.getCart(token)
//         ]);
//
//         // اگر هر دو درخواست موفقیت‌آمیز بود، داده‌ها را برمی‌گردانیم
//         return { user: userResponse.user, cart: cartResponse };
//
//     } catch (error) {
//         // اگر هر کدام از درخواست‌ها با خطا مواجه شود (مثلاً توکن نامعتبر)،
//         // فرض می‌کنیم کاربر لاگین نکرده است.
//         console.error("Failed to fetch initial data:", error);
//         return { user: null, cart: null };
//     }
// }

export const metadata: Metadata = {
    title: "فروشگاه بلک شاپ",
    description: "یک فروشگاه مدرن ساخته شده با Next.js و Go",
};

// RootLayout یک Server Component است و به صورت async اجرا می‌شود
export default async function RootLayout({
                                             children,
                                         }: Readonly<{
    children: React.ReactNode;
}>) {
    // اطلاعات اولیه را قبل از رندر شدن صفحه دریافت می‌کنیم
    // const { user, cart } = await getInitialData();

    // تعداد کل آیتم‌های سبد خرید را محاسبه می‌کنیم
    // const totalCartQuantity = cart?.items?.reduce((total: number, item: { quantity: number }) => {
    //     return total + item.quantity;
    // }, 0) || 0;

    const user = null;
    const totalCartQuantity = 0;

    return (
        <html lang="fa" dir="rtl">
        <body className="bg-gray-50 text-gray-800">
        <Header user={user} cartItemCount={totalCartQuantity} />
        <main className="container mx-auto p-4 sm:p-6 lg:p-8">
            {children}
        </main>
        </body>
        </html>
    );
}
