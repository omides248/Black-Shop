"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function logout() {
    let cookie = await cookies()
    cookie.delete('session');
    redirect('/auth/login');
}