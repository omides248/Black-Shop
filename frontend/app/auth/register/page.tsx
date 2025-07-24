"use client";

import { useActionState } from "react";
import { registerUser, FormState } from "./actions";
import { SubmitButton } from "./submit-button";

const initialState: FormState = {
    error: null,
    success: false,
};

export default function RegisterPage() {
    const [state, formAction] = useActionState(registerUser, initialState);

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-8 bg-gray-50">
            <div className="w-full max-w-md bg-white p-8 rounded-lg shadow-md">
                <h1 className="text-3xl font-bold mb-6 text-center text-gray-800">
                    ایجاد حساب کاربری
                </h1>
                <form action={formAction} className="space-y-6">
                    <div>
                        <label
                            htmlFor="name"
                            className="block text-sm font-medium text-gray-700"
                        >
                            نام
                        </label>
                        <input id="name" name="name" type="text" required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm" />
                    </div>
                    <div>
                        <label
                            htmlFor="email"
                            className="block text-sm font-medium text-gray-700"
                        >
                            ایمیل
                        </label>
                        <input id="email" name="email" type="email" required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm" />
                    </div>
                    <div>
                        <label
                            htmlFor="password"
                            className="block text-sm font-medium text-gray-700"
                        >
                            رمز عبور
                        </label>
                        <input id="password" name="password" type="password" required className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm" />
                    </div>

                    {state.error && (
                        <p className="text-sm text-red-500">{state.error}</p>
                    )}

                    {state.success && (
                        <p className="text-sm text-green-500">
                            ثبت‌نام با موفقیت انجام شد!
                        </p>
                    )}

                    <SubmitButton />
                </form>
            </div>
        </main>
    );
}