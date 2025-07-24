"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

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
        const res = await fetch("http://localhost:8081/v1/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        const data = await res.json();

        if (!res.ok) {
            return {
                error: data.message || "Failed to login.",
                success: false,
            };
        }

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
        console.error("Login error:", error);
        return { error: "An unexpected error occurred.", success: false };
    }

    redirect('/');
}
