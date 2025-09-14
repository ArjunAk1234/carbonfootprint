import { registerUser, loginUser } from './api.js';
import { createMessage, clearMessages } from './components.js';

document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const loginForm = document.getElementById('loginForm');
    const messageContainer = document.getElementById('messageContainer');

    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            clearMessages(messageContainer);

            const name = registerForm.name.value;
            const email = registerForm.email.value;
            const password = registerForm.password.value;
            const role = registerForm.role.value;

            try {
                const data = await registerUser({ name, email, password, role });
                messageContainer.appendChild(createMessage('success', data.message || 'Registration successful!'));
                registerForm.reset();
            } catch (error) {
                messageContainer.appendChild(createMessage('error', error.message || 'Registration failed.'));
            }
        });
    }

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            clearMessages(messageContainer);

            const email = loginForm.email.value;
            const password = loginForm.password.value;

            try {
                const data = await loginUser({ email, password });
                localStorage.setItem('jwt_token', data.token);
                localStorage.setItem('user_role', data.role); // Store user role
                messageContainer.appendChild(createMessage('success', data.message || 'Login successful! Redirecting...'));
                window.location.href = '/dashboard.html'; // Redirect to dashboard
            } catch (error) {
                messageContainer.appendChild(createMessage('error', error.message || 'Login failed.'));
            }
        });
    }
});