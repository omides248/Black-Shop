import Link from "next/link";

interface Product {
    id: string;
    name: string;
}

async function getProducts(): Promise<Product[]> {
    try {
        const res = await fetch("http://localhost:8080/v1/products", {
            cache: "no-cache",
        });

        if (!res.ok) {
            throw new Error("Failed to fetch the product list");
        }

        const data = await res.json();
        return data.products || [];
    } catch (error) {
        console.error('Error fetching products:', error);
        return []
    }
}

export default async function HomePage() {
    const products = await getProducts();

    return (
        <main className="flex min-h-screen flex-col items-center p-8 bg-gray-50">
            <h1 className="text-4xl font-bold mb-8 text-gray-800">فروشگاه ما</h1>

            {products.length > 0 ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 w-full max-w-7xl">
                    {products.map((product) => (
                        <Link href={`/products/${product.id}`} key={product.id}>
                            <div
                                className="bg-white rounded-lg shadow-md p-6 cursor-pointer transition-transform hover:scale-105">
                                <h2 className="text-xl font-semibold text-gray-700">{product.name}</h2>
                                <p className="text-sm text-gray-500 mt-2 font-mono break-all">{product.id}</p>
                            </div>
                        </Link>
                    ))}
                </div>
            ) : (
                <p className="text-xl text-gray-600">محصولی برای نمایش وجود ندارد.</p>
            )}
        </main>
    );
}
