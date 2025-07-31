"use client"; // This is our interactive client component

import { useActionState, useEffect, useState } from "react";
import { addItemToCart, CartActionState } from "@/app/cart/actions";
import { useFormStatus } from "react-dom";

const initialState: CartActionState = {};

function SubmitButton() {
    const { pending } = useFormStatus();
    return (
        <button
            type="submit"
            disabled={pending}
            className="w-full px-6 py-3 mt-4 text-white bg-indigo-600 rounded-md hover:bg-indigo-700 disabled:bg-indigo-400"
        >
            {pending ? "در حال افزودن..." : "افزودن به سبد خرید"}
        </button>
    );
}

export default function AddToCartButton({ productId }: { productId: string }) {
    const [state, formAction] = useActionState(addItemToCart, initialState);
    const [message, setMessage] = useState<string | null>(null);

    useEffect(() => {
        if (state.success) {
            setMessage("محصول با موفقیت به سبد خرید اضافه شد!");
            const timer = setTimeout(() => setMessage(null), 3000);
            return () => clearTimeout(timer);
        }
        if (state.error) {
            setMessage(state.error);
        }
    }, [state]);

    return (
        <form action={formAction}>
            <input type="hidden" name="productId" value={productId} />
            <input type="hidden" name="quantity" value="1" />
            <SubmitButton />
            {message && (
                <p className={`mt-2 text-sm ${state.error ? 'text-red-500' : 'text-green-600'}`}>
                    {message}
                </p>
            )}
        </form>
    );
}