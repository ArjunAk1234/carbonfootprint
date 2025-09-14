const API_BASE_URL = 'http://localhost:8080'; // Your Go backend URL

async function request(method, endpoint, data = null, requiresAuth = true) {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = {
        'Content-Type': 'application/json',
    };

    if (requiresAuth) {
        const token = localStorage.getItem('jwt_token');
        if (!token) {
            console.error('No JWT token found. Redirecting to login.');
            window.location.href = '/index.html'; // Redirect to login if token is missing
            return null; // Or throw an error
        }
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config = {
        method: method,
        headers: headers,
    };

    if (data) {
        config.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(url, config);
        const responseData = await response.json();

        if (!response.ok) {
            console.error(`API Error: ${response.status} - ${JSON.stringify(responseData)}`);
            // Specific handling for 401/403 errors
            if (response.status === 401 || response.status === 403) {
                 // Potentially refresh token or redirect to login
                 alert('Authentication expired or insufficient permissions. Please log in again.');
                 localStorage.removeItem('jwt_token');
                 window.location.href = '/index.html';
            }
            throw new Error(responseData.error || 'Something went wrong');
        }
        return responseData;

    } catch (error) {
        console.error('Fetch error:', error);
        throw error;
    }
}

// --- Specific API Functions ---

// Auth
export const registerUser = (userData) => request('POST', '/auth/register', userData, false);
export const loginUser = (credentials) => request('POST', '/auth/login', credentials, false);

// Users (Admin only)
export const getUsers = () => request('GET', '/users');
export const addUser = (userData) => request('POST', '/users', userData);
export const updateUser = (id, userData) => request('PUT', `/users/${id}`, userData);
export const deleteUser = (id) => request('DELETE', `/users/${id}`);

// Electric Consumption
export const getElectricConsumptions = () => request('GET', '/electric');
export const addElectricConsumption = (data) => request('POST', '/electric', data);
export const updateElectricConsumption = (id, data) => request('PUT', `/electric/${id}`, data);
export const deleteElectricConsumption = (id) => request('DELETE', `/electric/${id}`);

// Population
export const getPopulations = () => request('GET', '/population');
export const addPopulation = (data) => request('POST', '/population', data);
export const updatePopulation = (id, data) => request('PUT', `/population/${id}`, data);
export const deletePopulation = (id) => request('DELETE', `/population/${id}`);

// Transport
export const getTransports = () => request('GET', '/transport');
export const addTransport = (data) => request('POST', '/transport', data);
export const updateTransport = (id, data) => request('PUT', `/transport/${id}`, data);
export const deleteTransport = (id) => request('DELETE', `/transport/${id}`);

// Water Consumption
export const getWaterConsumptions = () => request('GET', '/water');
export const addWaterConsumption = (data) => request('POST', '/water', data);
export const updateWaterConsumption = (id, data) => request('PUT', `/water/${id}`, data);
export const deleteWaterConsumption = (id) => request('DELETE', `/water/${id}`);

// Waste
export const getWasteEntries = () => request('GET', '/waste');
export const addWasteEntry = (data) => request('POST', '/waste', data);
export const updateWasteEntry = (id, data) => request('PUT', `/waste/${id}`, data);
export const deleteWasteEntry = (id) => request('DELETE', `/waste/${id}`);

// Accommodation
export const getAccommodations = () => request('GET', '/accommodation');
export const addAccommodation = (data) => request('POST', '/accommodation', data);
export const updateAccommodation = (id, data) => request('PUT', `/accommodation/${id}`, data);
export const deleteAccommodation = (id) => request('DELETE', `/accommodation/${id}`);

// Goods Purchased
export const getGoodsPurchased = () => request('GET', '/goods');
export const addGoodsPurchased = (data) => request('POST', '/goods', data);
export const updateGoodsPurchased = (id, data) => request('PUT', `/goods/${id}`, data);
export const deleteGoodsPurchased = (id) => request('DELETE', `/goods/${id}`);

// Dashboard
export const getDashboardSummary = () => request('GET', '/dashboard');