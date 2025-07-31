// file: frontend/nextjs-app/components/Header.tsx
import Link from 'next/link';
import { logout } from '@/app/actions'; // Server Action برای خروج را ایمپورت می‌کنیم

// تعریف نوع داده برای پراپ‌های ورودی کامپوننت
interface HeaderProps {
    user: {
        name: string;
        // می‌توانید فیلدهای دیگر مثل id و email را هم اضافه کنید اگر نیاز بود
    } | null;
    cartItemCount: number;
}

export default function Header({ user, cartItemCount }: HeaderProps) {
    return (
        <header className="bg-white shadow-sm sticky top-0 z-50 border-b border-gray-200">
            <nav className="container mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    {/* بخش لوگو و لینک‌های اصلی */}
                    <div className="flex items-center gap-8">
                        <Link href="/">
                            <h1 className="text-2xl font-bold text-indigo-600 hover:text-indigo-800 transition-colors">
                                فروشگاه بلک شاپ
                            </h1>
                        </Link>
                        <div className="hidden md:flex md:gap-8">
                            <Link href="/" className="text-gray-600 hover:text-indigo-600 transition-colors">
                                خانه
                            </Link>
                            {/* لینک‌های دیگر مثل "دسته‌بندی‌ها" را می‌توان اینجا اضافه کرد */}
                        </div>
                    </div>

                    {/* بخش سبد خرید و پروفایل کاربر */}
                    <div className="flex items-center gap-6 text-sm font-medium">
                        <Link href="/cart" className="relative text-gray-600 hover:text-indigo-600 transition-colors p-2">
                            {/* آیکون سبد خرید */}
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                            </svg>
                            {cartItemCount > 0 && (
                                <span className="absolute top-0 right-0 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center transform translate-x-1/2 -translate-y-1/2">
                        {cartItemCount}
                    </span>
                            )}
                        </Link>

                        <div className="w-px h-6 bg-gray-200"></div> {/* جداکننده */}

                        {user ? (
                            // حالت لاگین کرده
                            <div className="flex items-center gap-4">
                                <Link href="/profile" className="text-gray-700 hover:text-indigo-600 transition-colors">
                                    خوش آمدید، <span className="font-bold">{user.name}</span>
                                </Link>
                                {/* فرم خروج که Server Action را صدا می‌زند */}
                                <form action={logout}>
                                    <button
                                        type="submit"
                                        className="text-red-500 hover:text-red-700 transition-colors"
                                    >
                                        خروج
                                    </button>
                                </form>
                            </div>
                        ) : (
                            // حالت لاگین نکرده
                            <div className="flex items-center gap-4">
                                <Link href="/auth/login" className="text-gray-700 hover:text-indigo-600 transition-colors">
                                    ورود
                                </Link>
                                <Link
                                    href="/auth/register"
                                    className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors"
                                >
                                    ثبت‌نام
                                </Link>
                            </div>
                        )}
                    </div>
                </div>
            </nav>
        </header>
    );
}
