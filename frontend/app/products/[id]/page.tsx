// file: frontend/nextjs-app/app/products/[id]/page.tsx
import { catalogAPI } from '@/lib/api/client';
import AddToCartForm from './AddToCartForm';

export const dynamic = 'force-dynamic';

async function getProduct(id: string) {
    try {
        const data = await catalogAPI.getProduct(id);
        return data;
    } catch (error) {
        console.error("Error fetching product:", error);
        return null;
    }
}

export default async function ProductPage({ params }: { params: { id: string } }) {
    const awaitedParams = await params;
    const id = awaitedParams.id;
    const product = await getProduct(id);

    if (!product) {
        return (
            <main className="flex min-h-screen flex-col items-center justify-center p-24">
                <h1 className="text-2xl font-bold">محصول یافت نشد!</h1>
            </main>
        );
    }

    return (
        <main className="container mx-auto p-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
                {/* Product Image Placeholder */}
                <div className="bg-gray-100 p-8 rounded-lg flex items-center justify-center aspect-square">
                    <span className="text-gray-400 text-xl">تصویر محصول</span>
                </div>

                {/* Product Details */}
                <div className="flex flex-col justify-center">
                    <h1 className="text-4xl font-bold text-gray-800">{product.name}</h1>

                    <div className="text-lg text-gray-500 font-mono mt-2 bg-gray-50 p-2 rounded-md inline-block self-start">
                        ID: {product.id}
                    </div>

                    <div className="mt-8 border-t pt-6">
                        <AddToCartForm productId={product.id} />
                    </div>
                </div>
            </div>
        </main>
    );
}
