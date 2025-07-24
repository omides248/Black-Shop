async function getProduct(id: string) {
    try {
        const res = await fetch(`http://127.0.0.1:8080/v1/products/${id}`, {
            cache: 'no-store',
        });

        if (!res.ok) {
            throw new Error('Failed to fetch product');
        }

        return await res.json();
    } catch (error) {
        console.error("Error fetching product:", error);
        return null;
    }
}

export default async function ProductPage({ params }: { params: any }) {
    const id = (await params)?.id;

    const product = await getProduct(id);

    if (!product) {
        return (
            <main className="flex min-h-screen flex-col items-center justify-center p-24">
                <h1 className="text-2xl font-bold">محصول یافت نشد!</h1>
            </main>
        );
    }

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-gray-100">
            <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
                <h1 className="text-2xl font-bold mb-4 text-gray-800">جزئیات محصول</h1>
                <div className="text-lg space-y-2">
                    <p>
                        <span className="font-semibold">شناسه:</span>
                        <span className="ml-2 font-mono p-1 bg-gray-200 rounded">{product.id}</span>
                    </p>
                    <p>
                        <span className="font-semibold">نام:</span>
                        <span className="ml-2">{product.name}</span>
                    </p>
                </div>
            </div>
        </main>
    );
}
