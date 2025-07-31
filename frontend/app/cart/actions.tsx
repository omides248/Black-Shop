"use server";

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache'; // ✅ ۱. ایمپورت کردن تابع
import { orderAPI } from '@/lib/api/client';

export interface CartActionState {
    error?: string;
    success?: boolean;
}

export async function addItemToCart(
    previousState: CartActionState,
    formData: FormData
): Promise<CartActionState> {
    const productId = formData.get("productId") as string;
    const quantity = parseInt(formData.get("quantity") as string, 10);


    const cookieStore = await cookies()
    const token = cookieStore.get('session')?.value;
    if (!token) {
        return { error: 'Please log in to add items to your cart.' };
    }

    try {
        await orderAPI.addItem({ productId, quantity }, token);

        // ✅ ۲. به Next.js می‌گوییم که layout اصلی را دوباره رندر کند
        // چون هدر در آن قرار دارد، با این کار تعداد آیتم‌ها آپدیت می‌شود
        revalidatePath('/');

        return { success: true };
    } catch (error: any) {
        console.error("Add to cart error:", error);
        return { error: error.message || "Could not add item to cart." };
    }
}



