import Link from 'next/link';
import { logout } from '@/app/actions';

interface HeaderProps {
    user: {
        id: string;
        name: string;
        email: string;
    } | null;
}

export default function Header({ user }: HeaderProps) {
    return (
        <header className="bg-white shadow-sm sticky top-0 z-50">
            <nav className="container mx-auto px-6 py-4 flex justify-between items-center">
                <Link href="/">
                    <h1 className="text-2xl font-bold text-indigo-600 hover:text-indigo-800 transition-colors">
                        فروشگاه بلک شاپ
                    </h1>
                </Link>
                <div className="flex items-center gap-6 text-sm font-medium">
                    {user ? (
                        <div className="flex items-center gap-4">
              <span className="text-gray-600">
                خوش آمدید، <span className="font-bold">{user.name}</span>
              </span>
                            <Link href="/profile" className="text-gray-700 hover:text-indigo-600 transition-colors">
                                پروفایل
                            </Link>
                            <form action={logout}>
                                <button
                                    type="submit"
                                    className="bg-red-500 text-white px-3 py-1.5 rounded-md hover:bg-red-600 transition-colors text-xs font-semibold"
                                >
                                    خروج
                                </button>
                            </form>
                        </div>
                    ) : (
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
            </nav>
        </header>
    );
}