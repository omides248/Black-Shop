import AddToCartButton from './AddToCartButton'; // We will create this next

export default function AddToCartForm({ productId }: { productId: string }) {
    return (
        <div>
            <AddToCartButton productId={productId} />
        </div>
    );
}