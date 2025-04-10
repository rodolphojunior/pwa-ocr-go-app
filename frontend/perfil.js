<script>

  document.addEventListener("DOMContentLoaded", () => {
  const token = localStorage.getItem("jwt");

      if (!token) {
        alert("⚠️ Sessão expirada. Faça login novamente.");
        window.location.href = "/";
        return;
      }
    });

//document.addEventListener("DOMContentLoaded", () => {
//  const token = localStorage.getItem("jwt");
//
//  if (!token) {
//    alert("Sessão expirada. Faça login novamente.");
//    window.location.href = "/";
//    return;
//  }

  // Exibe mensagem de boas-vindas
  document.getElementById("status").textContent = "🔐 Acesso autorizado com JWT";

  // Exemplo: buscar dados do perfil do usuário (ajuste conforme seu backend)
  fetch("/api/perfil", {
    headers: {
      Authorization: "Bearer " + token
    }
  })
    .then(res => {
      if (!res.ok) throw new Error("Não foi possível carregar perfil.");
      return res.json();
    })
    .then(data => {
      // Aqui você pode preencher os campos da tela com os dados do usuário
      console.log("📄 Dados do perfil:", data);
    })
    .catch(err => {
      console.error("Erro ao carregar perfil:", err);
      document.getElementById("status").textContent = "❌ Erro ao carregar perfil.";
    });

  // Service Worker (mover para arquivo JS, evitando inline no HTML)
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/service-worker.js')
      .then(reg => console.log('✅ Service Worker registrado:', reg))
      .catch(err => console.warn('⚠️ Falha ao registrar Service Worker:', err));
  }
});


  const CHAVES_KEY = "user_api_keys";
  const lista = document.getElementById("lista-chaves");

  function carregarChaves() {
    lista.innerHTML = "";
    const salvas = JSON.parse(localStorage.getItem(CHAVES_KEY) || "[]");
    salvas.forEach((chave, idx) => {
      const div = document.createElement("div");
      div.className = "api-entry";
      div.innerHTML = `
        <input type="text" value="${chave.nome}" disabled>
        <input type="text" value="${chave.valor}" disabled>
        <button onclick="removerChave(${idx})">❌</button>
      `;
      lista.appendChild(div);
    });
  }

  function removerChave(idx) {
    const chaves = JSON.parse(localStorage.getItem(CHAVES_KEY) || "[]");
    chaves.splice(idx, 1);
    localStorage.setItem(CHAVES_KEY, JSON.stringify(chaves));
    carregarChaves();
  }

  document.getElementById("nova-chave-form").addEventListener("submit", e => {
    e.preventDefault();
    const nome = document.getElementById("nome-chave").value.trim();
    const valor = document.getElementById("valor-chave").value.trim();
    if (!nome || !valor) return alert("Preencha ambos os campos");
    const chaves = JSON.parse(localStorage.getItem(CHAVES_KEY) || "[]");
    chaves.push({ nome, valor });
    localStorage.setItem(CHAVES_KEY, JSON.stringify(chaves));
    e.target.reset();
    carregarChaves();
  });

  document.getElementById("senha-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const atual = document.getElementById("senha-atual").value;
    const nova = document.getElementById("nova-senha").value;
    const conf = document.getElementById("confirmar-senha").value;
    if (nova !== conf) return alert("As senhas não coincidem");

    const token = localStorage.getItem("jwt");
    const res = await fetch("/trocar-senha", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({ senha_atual: atual, nova_senha: nova })
    });

    if (res.ok) {
      alert("Senha atualizada com sucesso!");
      e.target.reset();
    } else {
      const err = await res.text();
      alert("Erro ao atualizar senha: " + err);
    }
  });

  carregarChaves();

    if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/service-worker.js')
    .then(reg => console.log('✅ Service Worker registrado:', reg))
    .catch(err => console.warn('⚠️ Falha no registro do Service Worker:', err));
    }


</script>


