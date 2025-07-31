const API_PROTOCOL = process.env.NEXT_PUBLIC_API_PROTOCOL || 'http';


function getBaseUrl(service: 'catalog' | 'identity' | 'order'): string {
    let baseUrl: string | undefined;

    switch (service) {
        case 'catalog':
            baseUrl = process.env.NEXT_PUBLIC_CATALOG_API_BASE_URL;
            if (!baseUrl) throw new Error("Catalog API URL is not configured.");
            break;
        case 'identity':
            baseUrl = process.env.NEXT_PUBLIC_IDENTITY_API_BASE_URL;
            if (!baseUrl) throw new Error("Identity API URL is not configured.");
            break;
        case 'order':
            baseUrl = process.env.NEXT_PUBLIC_ORDER_API_BASE_URL;
            if (!baseUrl) throw new Error("Order API URL is not configured.");
            break;
        default:
            throw new Error("Invalid service name provided.");
    }

    return `${API_PROTOCOL}://${baseUrl}/v1`;
}


async function apiFetch(url: string, options: RequestInit = {}) {
    try {
        const res = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            cache: 'no-store',
        });

        if (!res.ok) {
            const errorData = await res.json();

            const apiError = {
                message: errorData.message || 'An API error occurred',
                code: errorData.code || 0,
                details: errorData.details || [],
                status: res.status,
            };
            throw apiError;
        }

        const text = await res.text();
        return text ? JSON.parse(text) : null;

    } catch (error: any) {
        console.error(`API fetch error for ${url}:`, error);
        // اگر خطای شبکه یا خطای دیگری قبل از دریافت پاسخ JSON رخ دهد
        if (error.message && !error.code) { // اگر یک Error استاندارد JS باشد
            throw { message: error.message, code: 0, details: [], status: 0 };
        }
        throw error; // خطای ساختاریافته را دوباره throw می‌کنیم
    }
}

export const catalogAPI = {
    getProduct: (id: string) => apiFetch(`${getBaseUrl('catalog')}/products/${id}`),
    getProducts: () => apiFetch(`${getBaseUrl('catalog')}/products`),
    createProduct: (data: { name: string }) => apiFetch(`${getBaseUrl('catalog')}/products`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
    listCategories: () => apiFetch(`${getBaseUrl('catalog')}/categories`),
    createCategory: (data: { name: string; imageUrl?: string | null; parentId?: string | null }) => apiFetch(`${getBaseUrl('catalog')}/categories`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
};

export const identityAPI = {
    register: (data: any) => apiFetch(`${getBaseUrl('identity')}/auth/register`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
    login: (data: any) => apiFetch(`${getBaseUrl('identity')}/auth/login`, {
        method: 'POST',
        body: JSON.stringify(data),
    }),
    getProfile: (token: string) => apiFetch(`${getBaseUrl('identity')}/users/me`, {
        headers: {'Authorization': `Bearer ${token}`},
    }),
};


export const orderAPI = {
    addItem: (data: { productId: string; quantity: number }, token: string) =>
        apiFetch(`${getBaseUrl('order')}/cart/items`, {
            method: 'POST',
            headers: {'Authorization': `Bearer ${token}`},
            body: JSON.stringify(data),
        }),
    getCart: (token: string) => apiFetch(`${getBaseUrl('order')}/cart`, {
        headers: { 'Authorization': `Bearer ${token}` },
    }),
};
