import {cookies} from 'next/headers';
import {identityAPI} from '@/lib/api/client';


async function getProfile() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session')?.value;

    if (!token) {
        return null;
    }

    try {
        const data = await identityAPI.getProfile(token);
        return data.user;
    } catch (error) {
        console.error("Failed to fetch profile:", error);
        return null;
    }
}

export default async function ProfilePage() {
    const user = await getProfile();

    return (
        <div className="container mx-auto p-8 max-w-2xl">
            <h1 className="text-3xl font-bold text-gray-800 border-b pb-4 mb-6">
                پروفایل کاربری
            </h1>
            {user ? (
                <div className="bg-white p-6 rounded-lg shadow-sm space-y-4 text-lg">
                    <div>
                        <span className="font-semibold text-gray-600">شناسه کاربر:</span>
                        <p className="text-gray-800 font-mono text-sm mt-1">{user.id}</p>
                    </div>
                    <div>
                        <span className="font-semibold text-gray-600">نام:</span>
                        <p className="text-gray-800 mt-1">{user.name}</p>
                    </div>
                    <div>
                        <span className="font-semibold text-gray-600">ایمیل:</span>
                        <p className="text-gray-800 mt-1">{user.email}</p>
                    </div>
                </div>
            ) : (
                <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md">
                    <p className="font-bold">خطا</p>
                    <p>اطلاعات پروفایل شما قابل دریافت نیست. لطفاً دوباره وارد شوید.</p>
                </div>
            )}
        </div>
    );
}