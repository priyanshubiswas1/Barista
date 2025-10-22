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
    const response = await api.get(`/sessions?shopid=${shopid}&tableid=${tableid}`);
    return response.status === 200 ? true : false
}

export const verifyUser = async () => {
    const response = await api.get(`/auth/verify`);
    return response.status === 200 ? true : false
}
