export function createMessage(type, text) {
    const div = document.createElement('div');
    div.className = `message ${type}`;
    div.textContent = text;
    return div;
}

export function clearMessages(container) {
    container.querySelectorAll('.message').forEach(msg => msg.remove());
}

export function createTable(headers, data, idKey = 'id') {
    const table = document.createElement('table');
    const thead = table.createTHead();
    const tbody = table.createTBody();

    const headerRow = thead.insertRow();
    headers.forEach(headerText => {
        const th = document.createElement('th');
        th.textContent = headerText;
        headerRow.appendChild(th);
    });
    // Add a column for actions if needed
    // const actionsTh = document.createElement('th');
    // actionsTh.textContent = 'Actions';
    // headerRow.appendChild(actionsTh);


    data.forEach(item => {
        const row = tbody.insertRow();
        headers.forEach(headerText => {
            const td = row.insertCell();
            const key = headerText.toLowerCase().replace(/ /g, '_'); // Convert "Item Name" to "item_name"
            let value = item[key];
            if (value instanceof Object && value.constructor === Date) {
                 value = value.toISOString().split('T')[0]; // Format Date objects
            } else if (typeof value === 'number' && !Number.isInteger(value)) {
                 value = value.toFixed(2); // Format floats
            }
            td.textContent = value || '-';
        });
        // Add action buttons if needed
        // const actionsTd = row.insertCell();
        // const editBtn = document.createElement('button');
        // editBtn.textContent = 'Edit';
        // editBtn.onclick = () => console.log('Edit', item[idKey]);
        // const deleteBtn = document.createElement('button');
        // deleteBtn.textContent = 'Delete';
        // deleteBtn.onclick = () => console.log('Delete', item[idKey]);
        // actionsTd.appendChild(editBtn);
        // actionsTd.appendChild(deleteBtn);
    });

    return table;
}