import * as api from './api.js';
import { createMessage, clearMessages, createTable } from './components.js';

const messageContainer = document.getElementById('messageContainer');
const dashboardSummarySection = document.getElementById('dashboardSummary');
const electricDataList = document.getElementById('electricDataList');
const populationDataList = document.getElementById('populationDataList');
const transportDataList = document.getElementById('transportDataList');
const waterDataList = document.getElementById('waterDataList');
const wasteDataList = document.getElementById('wasteDataList');
const accommodationDataList = document.getElementById('accommodationDataList');
const goodsDataList = document.getElementById('goodsDataList');

// Data Entry Forms
const addElectricForm = document.getElementById('addElectricForm');
const addPopulationForm = document.getElementById('addPopulationForm');
const addTransportForm = document.getElementById('addTransportForm');
const addWaterForm = document.getElementById('addWaterForm');
const addWasteForm = document.getElementById('addWasteForm');
const addAccommodationForm = document.getElementById('addAccommodationForm');
const addGoodsForm = document.getElementById('addGoodsForm');

document.addEventListener('DOMContentLoaded', async () => {
    // Check for token, redirect if not present
    if (!localStorage.getItem('jwt_token')) {
        window.location.href = '/index.html';
        return;
    }

    // Handle logout
    document.getElementById('logoutButton').addEventListener('click', () => {
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('user_role');
        window.location.href = '/index.html';
    });

    // --- Load Data on Dashboard ---
    await loadDashboardSummary();
    await loadElectricConsumptions();
    await loadPopulations();
    await loadTransports();
    await loadWaterConsumptions();
    await loadWasteEntries();
    await loadAccommodations();
    await loadGoodsPurchased();

    // --- Setup Form Submissions ---
    if (addElectricForm) addElectricForm.addEventListener('submit', handleAddElectric);
    if (addPopulationForm) addPopulationForm.addEventListener('submit', handleAddPopulation);
    if (addTransportForm) addTransportForm.addEventListener('submit', handleAddTransport);
    if (addWaterForm) addWaterForm.addEventListener('submit', handleAddWater);
    if (addWasteForm) addWasteForm.addEventListener('submit', handleAddWaste);
    if (addAccommodationForm) addAccommodationForm.addEventListener('submit', handleAddAccommodation);
    if (addGoodsForm) addGoodsForm.addEventListener('submit', handleAddGoods);
});

async function loadDashboardSummary() {
    try {
        const summary = await api.getDashboardSummary();
        dashboardSummarySection.innerHTML = `
            <h2>Overall Carbon Footprint</h2>
            <div class="dashboard-grid">
                <div class="dashboard-card">
                    <h3>Total CO2e</h3>
                    <p>${summary.total_carbon_footprint_co2e.toFixed(2)} kg CO2e</p>
                </div>
                <div class="dashboard-card">
                    <h3>Total Population</h3>
                    <p>${summary.total_population}</p>
                </div>
                <div class="dashboard-card">
                    <h3>Per-Capita CO2e</h3>
                    <p>${summary.per_capita_footprint_co2e.toFixed(2)} kg CO2e</p>
                </div>
            </div>
            <h3>Component Breakdown (kg CO2e)</h3>
            <ul>
                ${Object.entries(summary.component_breakdown).map(([key, value]) => `
                    <li><strong>${key}:</strong> ${value.toFixed(2)}</li>
                `).join('')}
            </ul>
            <!-- Trends would go here, requiring more complex chart libraries -->
        `;
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load dashboard summary: ${error.message}`));
    }
}

async function loadElectricConsumptions() {
    try {
        const data = await api.getElectricConsumptions();
        if (data.length > 0) {
            const headers = ['ID', 'Source', 'KWH', 'Fuel Liters', 'Hours', 'Date', 'Location'];
            electricDataList.innerHTML = ''; // Clear previous data
            electricDataList.appendChild(createTable(headers, data));
        } else {
            electricDataList.innerHTML = '<p>No electric consumption data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load electric consumption: ${error.message}`));
    }
}

async function loadPopulations() {
    try {
        const data = await api.getPopulations();
        if (data.length > 0) {
            const headers = ['ID', 'Registered Count', 'Floating Count', 'Date', 'Location'];
            populationDataList.innerHTML = '';
            populationDataList.appendChild(createTable(headers, data));
        } else {
            populationDataList.innerHTML = '<p>No population data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load population data: ${error.message}`));
    }
}

async function loadTransports() {
    try {
        const data = await api.getTransports();
        if (data.length > 0) {
            const headers = ['ID', 'Vehicle Type', 'Fuel Type', 'Distance KM', 'Fuel Liters', 'Date', 'Location'];
            transportDataList.innerHTML = '';
            transportDataList.appendChild(createTable(headers, data));
        } else {
            transportDataList.innerHTML = '<p>No transport data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load transport data: ${error.message}`));
    }
}

async function loadWaterConsumptions() {
    try {
        const data = await api.getWaterConsumptions();
        if (data.length > 0) {
            const headers = ['ID', 'Meter Reading', 'Date', 'Location'];
            waterDataList.innerHTML = '';
            waterDataList.appendChild(createTable(headers, data));
        } else {
            waterDataList.innerHTML = '<p>No water consumption data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load water consumption: ${error.message}`));
    }
}

async function loadWasteEntries() {
    try {
        const data = await api.getWasteEntries();
        if (data.length > 0) {
            const headers = ['ID', 'Spot Name', 'Waste Type', 'Weight KG', 'Date', 'Location'];
            wasteDataList.innerHTML = '';
            wasteDataList.appendChild(createTable(headers, data));
        } else {
            wasteDataList.innerHTML = '<p>No waste data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load waste data: ${error.message}`));
    }
}

async function loadAccommodations() {
    try {
        const data = await api.getAccommodations();
        if (data.length > 0) {
            const headers = ['ID', 'People Count', 'Nights', 'Date', 'Location'];
            accommodationDataList.innerHTML = '';
            accommodationDataList.appendChild(createTable(headers, data));
        } else {
            accommodationDataList.innerHTML = '<p>No accommodation data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load accommodation data: ${error.message}`));
    }
}

async function loadGoodsPurchased() {
    try {
        const data = await api.getGoodsPurchased();
        if (data.length > 0) {
            const headers = ['ID', 'Item Name', 'Quantity', 'Cost', 'Date', 'Location'];
            goodsDataList.innerHTML = '';
            goodsDataList.appendChild(createTable(headers, data));
        } else {
            goodsDataList.innerHTML = '<p>No goods purchased data found.</p>';
        }
    } catch (error) {
        messageContainer.appendChild(createMessage('error', `Failed to load goods purchased data: ${error.message}`));
    }
}


// --- Handlers for Adding Data ---

async function handleAddElectric(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        source: form.source.value,
        kwh: parseFloat(form.kwh.value) || 0,
        fuel_liters: parseFloat(form.fuel_liters.value) || 0,
        hours: parseFloat(form.hours.value) || 0,
        date: form.date.value + "T00:00:00Z", // Convert to ISO string
        location: form.location.value || 'Overall',
    };

    try {
        await api.addElectricConsumption(data);
        messageContainer.appendChild(createMessage('success', 'Electric consumption added!'));
        form.reset();
        await loadElectricConsumptions();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add electric consumption.'));
    }
}

async function handleAddPopulation(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        registered_count: parseInt(form.registered_count.value),
        floating_count: parseInt(form.floating_count.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addPopulation(data);
        messageContainer.appendChild(createMessage('success', 'Population stats added!'));
        form.reset();
        await loadPopulations();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add population stats.'));
    }
}

async function handleAddTransport(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        vehicle_type: form.vehicle_type.value,
        fuel_type: form.fuel_type.value,
        distance_km: parseFloat(form.distance_km.value),
        fuel_liters: parseFloat(form.fuel_liters.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addTransport(data);
        messageContainer.appendChild(createMessage('success', 'Transport data added!'));
        form.reset();
        await loadTransports();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add transport data.'));
    }
}

async function handleAddWater(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        meter_reading: parseFloat(form.meter_reading.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addWaterConsumption(data);
        messageContainer.appendChild(createMessage('success', 'Water reading added!'));
        form.reset();
        await loadWaterConsumptions();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add water reading.'));
    }
}

async function handleAddWaste(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        spot_name: form.spot_name.value,
        waste_type: form.waste_type.value,
        weight_kg: parseFloat(form.weight_kg.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addWasteEntry(data);
        messageContainer.appendChild(createMessage('success', 'Waste entry added!'));
        form.reset();
        await loadWasteEntries();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add waste entry.'));
    }
}

async function handleAddAccommodation(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        people_count: parseInt(form.people_count.value),
        nights: parseInt(form.nights.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addAccommodation(data);
        messageContainer.appendChild(createMessage('success', 'Accommodation data added!'));
        form.reset();
        await loadAccommodations();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add accommodation data.'));
    }
}

async function handleAddGoods(e) {
    e.preventDefault();
    clearMessages(messageContainer);
    const form = e.target;
    const data = {
        item_name: form.item_name.value,
        quantity: parseInt(form.quantity.value),
        cost: parseFloat(form.cost.value),
        date: form.date.value + "T00:00:00Z",
        location: form.location.value || 'Overall',
    };

    try {
        await api.addGoodsPurchased(data);
        messageContainer.appendChild(createMessage('success', 'Goods purchased added!'));
        form.reset();
        await loadGoodsPurchased();
        await loadDashboardSummary();
    } catch (error) {
        messageContainer.appendChild(createMessage('error', error.message || 'Failed to add goods purchased.'));
    }
}