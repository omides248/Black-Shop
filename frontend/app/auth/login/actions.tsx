"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import {identityAPI} from "@/lib/api/client";

export interface FormState {
    error: string | null;
    success: boolean;
}

export async function loginUser(
    previousState: FormState,
    formData: FormData
): Promise<FormState> {
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
        const data = await identityAPI.login({ email, password });
        const token = data.token;
        if (token) {
            const cookieStore = await cookies();
            cookieStore.set("session", token, {
                httpOnly: true,
                secure: process.env.NODE_ENV === "production",
                path: "/",
                sameSite: "strict",
                maxAge: 60 * 60 * 24
            });
        } else {
            return { error: "Token not received from server.", success: false };
        }

    } catch (error) {
        console.error("Login action error:", error);
        return { error: "An unexpected error occurred.", success: false };
    }

    redirect('/');
}
