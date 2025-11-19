import axios from "axios"

const api = axios.create({
    baseURL: "http://localhost:8000/api",
    headers: {
        'Content-Type': 'application/json',
    },
    credentials: 'include'
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

export const registerUser = async (credentials) => {
    const response = await api.post(`/auth/onboard`, credentials);
    return response.status === 200 ? response.data : null
}

export const loginUser = async (credentials) => {
    const response = await api.post(`/auth/login`, credentials);
    return response.status === 200 ? response.data : null
}

export const verifyUser = async () => {
    const response = await api.get(`/auth/verify`);
    return response.status === 200 ? true : false
}

export const fetchTablesByShopId = async (shopId) => {
    const response = await api.get(`/tables/${shopId}`);
    return response.status === 200 ? response.data : null
};

export const createTableByShopId = async ({ shopId, tableNumber }) => {
    const response = await api.post(`/table`, {
        shopId: shopId,
        tableNumber: tableNumber
    });
    return response.status === 200 ? response.data : null
};

export const fetchCategoryItemsByShopId = async (id) => {
    const response = await api.get(`/category-items/${id}`);
    return response.status === 200 ? response.data : null
};

export const fetchCategoriesByShopId = async (id) => {
    const response = await api.get(`/category/${id}`);
    return response.status === 200 ? response.data : null
};

export const deleteCategoryByCategoryId = async ({ categoryId, shopId }) => {
    const response = await api.delete(`/category`, {
        data: { categoryId, shopId }
    });
    return response.status === 200 ? response.data : null
};

export const createCategory = async ({ name, shopId }) => {
    const response = await api.post(`/category`, { name: name, shopId: shopId });
    return response.status === 200 ? response.data : null
}

export const updateCategory = async ({ name, shopId, categoryId }) => {
    const response = await api.put(`/category`, { name: name, shopId: shopId, categoryId: categoryId });
    return response.status === 200 ? response.data : null
}

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}/admin`);
    return response.status === 200 ? response.data : null
};

export const addItem = async ({ name, shopId, categoryId, price, veg }) => {
    const response = await api.post(`/item`, { name: name, shopId: shopId, description: "test", categoryId: categoryId, veg: veg, price: price, image: "url", available: true });
    return response.status === 200 ? response.data : null
}

export const updateItem = async ({ itemId, shopId, itemName, categoryId, price, veg }) => {
    const response = await api.put(`/item`, { itemId: itemId, shopId: shopId, itemName: itemName, description: "test", categoryId: categoryId, veg: veg, price: price, image: "url", available: true });
    return response.status === 200 ? response.data : null
}

export const updateItemAvailable = async ({ itemId, shopId, available }) => {
    const response = await api.put(`/item-available`, { itemId: itemId, shopId: shopId, available: available });
    return response.status === 200 ? response.data : null
}

export const deleteItemById = async ({ itemId }) => {
    const response = await api.delete(`/item/${itemId}`);
    return response.status === 200 ? response.data : null
}

export const fetchActiveSessions = async (shopId) => {
    const response = await api.get(`/active-sessions?shopid=${shopId}`);
    return response.status === 200 ? response.data : null
};

export const fetchOrdersBySessionId = async (sessionId) => {
    const response = await api.get(`/orders/${sessionId}`);
    return response.status === 200 ? response.data : null
};

export const createOrder = async (payload) => {
    const response = await api.post(`/order-admin`, payload);
    return response.status === 200 ? response.data : null
};

export const getSessionSummary = async (payload) => {
    let url = `/history?shopid=${payload.shopId}&offset=${payload.offset}&limit=${payload.limit}`;
    if (payload.timeRange) {
        url += `&time-range=${payload.timeRange}`;
    }
    const response = await api.get(url);
    return response.status === 200 ? response.data : null
};