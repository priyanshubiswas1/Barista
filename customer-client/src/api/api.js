import axios from "axios"

const api = axios.create({
    baseURL: "http://localhost:8000/api",
    headers: {
        'Content-Type': 'application/json',
    },
    credentials: 'include'
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem("session-token");
    if (token) {
        config.headers.Session = `${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

export const sendOtp = async (payload) => {
    const response = await api.post(`/sessions/send-otp`, payload);
    return response.status === 200 ? response.data : null
}

export const createSession = async (payload) => {
    const response = await api.post(`/sessions`, payload);
    return response.status === 200 ? response.data : null
}

export const joinSession = async (payload) => {
    const response = await api.post(`/sessions/join`, payload);
    return response.status === 200 ? response.data : null
}

export const getSession = async ({ shopid, tableid }) => {
    try {
        const response = await api.get(`/sessions?shopid=${shopid}&tableid=${tableid}`);
        if (response.status === 200 && response.data) {
            if (response.data.sessionToken) {
                localStorage.setItem("session-token", response.data.sessionToken);
            }
            return response.data;
        }
        localStorage.setItem("session-token", "");
        return null;
    } catch (err) {
        localStorage.setItem("session-token", "");
        return null;
    }
}

export const placeOrder = async (payload) => {
    const response = await api.post(`/order-customer`, payload);
    if (response.status === 200 && response.data) {
        // Always update token - handles both new and existing sessions
        if (response.data.sessionToken) {
            localStorage.setItem("session-token", response.data.sessionToken);
        }
        return response.data;
    }
    return null;
}

export const getOrderByOrderId = async (id) => {
    const response = await api.get(`/order/${id}`);
    return response.status === 200 ? response.data : null
};

export const getOrdersBySessionId = async () => {
    const response = await api.get(`/orders`);
    return response.status === 200 ? response.data : null
};

export const verifyUser = async () => {
    const response = await api.get(`/auth/verify`);
    return response.status === 200 ? true : false
}

export const fetchItemsByShopId = async (id) => {
    const response = await api.get(`/items/${id}`);
    return response.status === 200 ? response.data : null
};

export const fetchTableByTableId = async ({ shopid, tableid }) => {
    const response = await api.get(`/customer/table/${tableid}`, {
        "shopId": shopid
    });
    return response.status === 200 ? response.data : null
};
