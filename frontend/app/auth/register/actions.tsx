"use server";

import { redirect } from 'next/navigation';
import { identityAPI } from '@/lib/api/client';

export interface FormState {
    error: string | null;
    success: boolean;
}


export async function registerUser(
    previousState: FormState,
    formData: FormData
): Promise<FormState> {

    const name = formData.get("name") as string;
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
        await identityAPI.register({ name, email, password });

    } catch (error: any) {
        console.error("Register action error:", error);

        return { error: error.message || "An unexpected error occurred.", success: false };
    }
    
    redirect('/auth/login');
}