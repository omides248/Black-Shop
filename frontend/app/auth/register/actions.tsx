"use server";

export interface FormState {
    error: string | null;
    success: boolean;
}

export async function registerUser(previousState: FormState, formData: FormData): Promise<FormState> {
    const name = formData.get("name") as string;
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
        const res = await fetch("http://localhost:8081/v1/auth/register", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name, email, password}),
        });

        if (!res.ok) {
            const errorData = await res.json();
            return {
                error: errorData.message || "Failed to register.",
                success: false,
            };
        }

        return {success: true, error: null};
    } catch (error) {
        console.error("Registration error:", error);
        return {error: "An unexpected error occurred.", success: false};
    }
}