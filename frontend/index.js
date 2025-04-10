 
      // Constantes e Elementos DOM
      const API_ENDPOINTS = {
        LOGIN: '/login',
        REGISTER: '/register'
      };
      const JWT_KEY = 'jwt';
      const DEFAULT_REDIRECT = '/img2txt';

      const forms = {
        login: document.getElementById('login-form'),
        register: document.getElementById('register-form'),
        status: document.getElementById('status'),
        logoutBox: document.getElementById('logout-box'),
        welcomeMsg: document.getElementById('welcome-msg')
      };

      // Funções Utilitárias
      const showFeedback = (message, type = 'error') => {
        forms.status.className = type;
        forms.status.textContent = message;
      };

      const setLoadingState = (isLoading) => {
        document.querySelectorAll('button[type="submit"]').forEach(btn => {
          btn.disabled = isLoading;
          btn.innerHTML = isLoading ? '⏳ Processando...' : btn.dataset.originalText;
        });
      };

      // Handlers de Formulário
    const handleLogin = async (data) => {
      try {
        setLoadingState(true);
        forms.status.textContent = '';

        const response = await fetch(API_ENDPOINTS.LOGIN, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(data)
        });

        if (!response.ok) {
          const error = await response.text();
          throw new Error(error || 'Erro no login');
        }

        const result = await response.json();
        localStorage.setItem(JWT_KEY, result.token); // ✅ Salva o token no localStorage
        console.log("🔑 Token salvo:", result.token);

        showFeedback('Login bem-sucedido! Redirecionando...', 'success');

        setTimeout(() => {
          window.location.href = DEFAULT_REDIRECT; // ✅ Redireciona para /img2txt
        }, 800);

      } catch (error) {
        showFeedback(`❌ Erro no login: ${error.message}`);
      } finally {
        setLoadingState(false);
      }
    };


      const handleRegister = async (data) => {
        try {
          const response = await fetch(API_ENDPOINTS.REGISTER, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
          });

          if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Erro no registro');
          }

          showFeedback('Registro realizado! Faça login.', 'success');
          document.getElementById('show-login').click();

        } catch (error) {
          showFeedback(`Falha no registro: ${error.message}`);
        }
      };

      // Event Listeners
      forms.login.onsubmit = async (e) => {
        e.preventDefault();
        setLoadingState(true);
        await handleLogin(Object.fromEntries(new FormData(forms.login)));
        setLoadingState(false);
      };

      forms.register.onsubmit = async (e) => {
        e.preventDefault();
        setLoadingState(true);
        await handleRegister(Object.fromEntries(new FormData(forms.register)));
        setLoadingState(false);
      };

      // Toggle entre Login/Registro
      document.getElementById('show-register').onclick = () => {
        forms.login.classList.add('hidden');
        forms.register.classList.remove('hidden');
        document.getElementById('title').textContent = 'Registro';
      };

      document.getElementById('show-login').onclick = () => {
        forms.register.classList.add('hidden');
        forms.login.classList.remove('hidden');
        document.getElementById('title').textContent = 'Login';
      };

      // Logout
      document.getElementById('logout-btn').onclick = () => {
        localStorage.removeItem(JWT_KEY);
        window.location.reload();
      };

      // Preservar textos dos botões
      document.querySelectorAll('button[type="submit"]').forEach(btn => {
        btn.dataset.originalText = btn.innerHTML;
      });

    // Service Worker
    if ('serviceWorker' in navigator) {
      window.addEventListener('load', () => {
        navigator.serviceWorker.register('/service-worker.js')
          .then(registration => {
            console.log('SW registrado:', registration);
          })
          .catch(error => {
            console.error('Falha no SW:', error);
          });
      });
    }


